package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/matsuyoshi30/gom2h"
	"gopkg.in/yaml.v2"
)

type JournalType int

const (
	TypeDaily JournalType = iota
	TypeWeekly
	TypeMonthly
)

type Config struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	TypeStr     string `yaml:"typestr"`
	Wd          string `yaml:"wd"`
	ContentDir  string `yaml:"contentDir"`
	TemplateDir string `yaml:"templateDir"`
	PublishDir  string `yaml:"publishDir"`
	Type        JournalType
	Posts       []Post
}

var config Config

type Post struct {
	Title     string
	Body      template.HTML
	Link      string
	UpdatedAt time.Time
}

func (config *Config) New(dirpath string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	wd = filepath.Join(wd, dirpath)

	dirs := []string{"public", "content", "static", "template"}
	for _, d := range dirs {
		if err := os.Mkdir(filepath.Join(wd, d), os.ModePerm); err != nil {
			return err
		}
	}

	files := []struct {
		dir      string
		filename string
		contents string
	}{
		{"template", "index.html.tmpl", indexTmpl},
		{"template", "content.md.tmpl", contentTmpl},
		{"", "config.yaml", configTmpl + "typestr: " + config.typeToString() + "\n" + "wd: " + wd + "\n"},
	}
	for _, f := range files {
		file := filepath.Join(wd, f.filename)
		if f.dir != "" {
			file = filepath.Join(wd, f.dir, f.filename)
		}

		if err := ioutil.WriteFile(file, []byte(f.contents), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (config *Config) typeToString() string {
	switch config.Type {
	case TypeDaily:
		return "TypeDaily"
	case TypeWeekly:
		return "TypeWeekly"
	case TypeMonthly:
		return "TypeMonthly"
	}

	return ""
}

func (config *Config) stringToType() JournalType {
	switch config.TypeStr {
	case "TypeDaily":
		return TypeDaily
	case "TypeWeekly":
		return TypeWeekly
	case "TypeMonthly":
		return TypeMonthly
	}

	return 0
}

func (config *Config) Post() error {
	year, month, day := time.Now().Date()

	// e.g.) 2020-05-24
	//   TypeDaily   -> `<contentDir>/2020/05/24.md`
	//   TypeWeekly  -> `<contentDir>/2020/05-18.md` check starting weekday
	//   TypeMonthly -> `<contentDir>/202005.md`
	contentDir := config.ContentDir
	if config.Type == TypeDaily {
		contentDir = filepath.Join(contentDir, strconv.Itoa(year), strconv.Itoa(int(month)))
	} else if config.Type == TypeWeekly {
		contentDir = filepath.Join(contentDir, strconv.Itoa(year))
	}
	if err := os.MkdirAll(contentDir, os.ModePerm); err != nil {
		return err
	}

	// read content template
	b, err := ioutil.ReadFile(filepath.Join(config.TemplateDir, "content.md.tmpl"))
	if err != nil {
		return err
	}

	// create post file
	var filename string
	switch config.Type {
	case TypeDaily:
		filename = strconv.Itoa(day) + ".md"
	case TypeWeekly:
		filename = checkWeekday(time.Now()) + ".md"
	case TypeMonthly:
		filename = strconv.Itoa(year) + strconv.Itoa(int(month)) + ".md"
	}
	if err := ioutil.WriteFile(filepath.Join(contentDir, filename), b, 0644); err != nil {
		return err
	}

	return nil
}

func checkWeekday(date time.Time) string {
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
	}
	_, month, day := date.Date()

	return strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
}

func (config *Config) Load(filename string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(filepath.Join(wd, filename))
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(b, config); err != nil {
		return err
	}

	if config.ContentDir == "" {
		config.ContentDir = filepath.Join(config.Wd, "content")
	}
	if config.TemplateDir == "" {
		config.TemplateDir = filepath.Join(config.Wd, "template")
	}
	if config.PublishDir == "" {
		config.PublishDir = filepath.Join(config.Wd, "public")
	}
	config.Type = config.stringToType()

	return nil
}

func (config *Config) Build(dest string) error {
	posts := make([]Post, 0)

	if err := filepath.Walk(config.ContentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		parentDir := filepath.Dir(path)
		if info.IsDir() {
			if path != config.ContentDir {
				if string(path[len(path)-5]) == "/" {
					if err := os.Mkdir(filepath.Join(dest, path[len(parentDir)+1:]), os.ModePerm); err != nil {
						return err
					}
				} else {
					y := parentDir[len(parentDir)-4:]
					m := path[len(parentDir)+1:]
					if err := os.Mkdir(filepath.Join(dest, y, m), os.ModePerm); err != nil {
						return err
					}
				}
			}
		} else {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			output, err := gom2h.Run(b)
			if err != nil {
				return err
			}

			filename := filepath.Base(path)
			ext := filepath.Ext(filename)
			title := filename[0 : len(filename)-len(ext)]
			htmlFile := title + ".html" // TypeMonthly

			if !strings.HasSuffix(parentDir, "content") {
				if string(parentDir[len(parentDir)-3]) != "/" { // TypeWeekly
					htmlFile = filepath.Join(parentDir[len(parentDir)-4:], htmlFile)
				} else { // TypeDaily
					m := parentDir[len(parentDir)-2:]
					y := parentDir[len(parentDir)-7 : len(parentDir)-3]
					htmlFile = filepath.Join(y, m, htmlFile)
				}
			}

			post := Post{
				Title:     title,
				Body:      template.HTML(output),
				Link:      "http://localhost:8080/" + htmlFile,
				UpdatedAt: info.ModTime(),
			}

			t, err := template.New("post").Parse(postTmpl)
			if err != nil {
				return err
			}
			f, err := os.Create(filepath.Join(dest, htmlFile))
			if err != nil {
				return err
			}
			if err := t.Execute(f, post); err != nil {
				return err
			}

			posts = append(posts, post)
		}
		return nil
	}); err != nil {
		return err
	}
	config.Posts = posts

	return nil
}

func generateTemplate(tmplName, tmplPath string) (*template.Template, error) {
	tmpl, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return nil, err
	}

	t, err := template.New(tmplName).Parse(string(tmpl))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (config *Config) createFileFromTemplate(tmpl *template.Template, dstPath string) error {
	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}

	return tmpl.Execute(f, config)
}

func (config *Config) Serve() error {
	dir, err := ioutil.TempDir("", "tmp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	// build html from markdown file
	if err := config.Build(dir); err != nil {
		return err
	}

	t, err := generateTemplate("index", filepath.Join(config.TemplateDir, "index.html.tmpl"))
	if err != nil {
		return err
	}
	if err := config.createFileFromTemplate(t, filepath.Join(dir, "index.html")); err != nil {
		return err
	}

	http.Handle("/", http.FileServer(http.Dir(dir)))
	return http.ListenAndServe(":8080", nil)
}

func (config *Config) Publish() error {
	if err := config.Build(config.PublishDir); err != nil {
		return err
	}

	t, err := generateTemplate("index", filepath.Join(config.TemplateDir, "index.html.tmpl"))
	if err != nil {
		return err
	}
	if err := config.createFileFromTemplate(t, filepath.Join(config.PublishDir, "index.html")); err != nil {
		return err
	}

	return nil
}
