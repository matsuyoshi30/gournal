package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
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
	StaticDir   string `yaml:"staticDir"`
	PublishDir  string `yaml:"publishDir"`
	Type        JournalType
	Posts       []Post
}

var config Config

type Post struct {
	Title      string
	Body       template.HTML
	PostYear   string
	PostMonth  string
	FromDate   string // use if TypeWeekly
	ToDate     string // use if TypeWeekly
	IsLastWeek bool   // use if TypeWeekly
	WeekNum    int    // use if TypeWeekly
	PostDate   string // use if TypeMonthly or TypeDaily
	BaseLink   string
	CSSPath    string
	Link       string
	UpdatedAt  time.Time
}

func createDirs(dir string) error {
	dirs := []string{"public", "content", "static", "template"}
	for _, d := range dirs {
		if err := os.Mkdir(filepath.Join(dir, d), os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func createFiles(dir string) error {
	files := []struct {
		dir      string
		filename string
		contents string
	}{
		{"template", "index.html.tmpl", indexTmpl},
		{"template", "post.html.tmpl", postTmpl},
		{"template", "content.md.tmpl", contentTmpl},
		{"static", "styles.css", cssTmpl},
		{"", "config.yaml", configTmpl + "typestr: " + config.typeToString() + "\n" + "wd: " + dir + "\n"},
	}
	for _, f := range files {
		file := filepath.Join(dir, f.filename)
		if f.dir != "" {
			file = filepath.Join(dir, f.dir, f.filename)
		}

		if err := ioutil.WriteFile(file, []byte(f.contents), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (config *Config) New(dirpath string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	projectDir := filepath.Join(wd, dirpath)
	if err := createDirs(projectDir); err != nil {
		return err
	}

	if err := createFiles(projectDir); err != nil {
		return err
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
	var contentDir, postFilename string
	switch config.Type {
	case TypeMonthly:
		contentDir = config.ContentDir
		postFilename = strconv.Itoa(year) + fmt.Sprintf("%02d", int(month)) + ".md"
	case TypeDaily:
		contentDir = filepath.Join(config.ContentDir, strconv.Itoa(year), fmt.Sprintf("%02d", int(month)))
		postFilename = strconv.Itoa(day) + ".md"
	case TypeWeekly:
		contentDir = filepath.Join(config.ContentDir, strconv.Itoa(year))
		postFilename = checkWeekday(time.Now()) + ".md"
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
	if err := ioutil.WriteFile(filepath.Join(contentDir, postFilename), b, 0644); err != nil {
		return err
	}

	return nil
}

func checkWeekday(date time.Time) string {
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
	}
	_, month, day := date.Date()

	return fmt.Sprintf("%02d", int(month)) + "-" + strconv.Itoa(day)
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

	config.Wd = wd
	if config.ContentDir == "" {
		config.ContentDir = filepath.Join(config.Wd, "content")
	}
	if config.TemplateDir == "" {
		config.TemplateDir = filepath.Join(config.Wd, "template")
	}
	if config.PublishDir == "" {
		config.PublishDir = filepath.Join(config.Wd, "public")
	}
	if config.StaticDir == "" {
		config.StaticDir = filepath.Join(config.Wd, "static")
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
				switch config.Type {
				case TypeWeekly:
					// <contentDir>/2020/
					if err := createIfNotExists(filepath.Join(dest, path[len(parentDir)+1:])); err != nil {
						return err
					}
				case TypeDaily:
					// <contentDir>/2020/05/
					y := parentDir[len(parentDir)-4:]
					m := path[len(parentDir)+1:]
					if err := createIfNotExists(filepath.Join(dest, y, m)); err != nil {
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

			post := Post{
				Title:      extractTitle(path),
				Body:       template.HTML(output),
				IsLastWeek: false,
				WeekNum:    0,
				BaseLink:   dest,
				CSSPath:    "./static/styles.css",
				UpdatedAt:  info.ModTime(),
			}

			var htmlFile string
			var yearStr, monthStr string
			switch config.Type {
			case TypeMonthly:
				// <contentDir>/202005.md
				yearStr = post.Title[:4]
				monthStr = post.Title[4:]
				htmlFile = post.Title + ".html"
			case TypeWeekly:
				// <contentDir>/2020/05-25.md
				yearStr = parentDir[len(parentDir)-4:]
				monthStr = post.Title[:2]
				post.FromDate = post.Title
				to, err := toWeekend(yearStr, post.Title)
				if err != nil {
					return err
				}
				post.ToDate = to
				w, err := time.Parse("2006-01-02", yearStr+"-"+post.Title)
				if err != nil {
					return err
				}
				_, post.WeekNum = w.ISOWeek()
				post.IsLastWeek = post.WeekNum-52 >= 0
				htmlFile = filepath.Join(yearStr, post.Title+".html")
				post.CSSPath = "../static/styles.css"
			case TypeDaily:
				// <contentDir>/2020/05/25.md
				yearStr = parentDir[len(parentDir)-7 : len(parentDir)-3]
				monthStr = parentDir[len(parentDir)-2:]
				post.PostDate = post.Title
				htmlFile = filepath.Join(yearStr, monthStr, post.Title+".html")
				post.CSSPath = "../../static/styles.css"
			}

			post.PostYear = yearStr
			post.PostMonth = monthStr
			post.Link = "./" + htmlFile

			t, err := generateTemplate("post", filepath.Join(config.TemplateDir, "post.html.tmpl"))
			if err != nil {
				return err
			}
			if err := config.createFileFromTemplate(t, filepath.Join(dest, htmlFile), post); err != nil {
				return err
			}

			posts = append(posts, post)
		}

		return nil
	}); err != nil {
		return err
	}
	config.Posts = reverse(posts)

	staticDestDir := filepath.Join(dest, "static")
	if err := createIfNotExists(staticDestDir); err != nil {
		return err
	}
	if err := copyDir(config.StaticDir, staticDestDir); err != nil {
		return err
	}

	return nil
}

func extractTitle(path string) string {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)

	return filename[0 : len(filename)-len(ext)]
}

func toWeekend(y, md string) (string, error) {
	date, err := time.Parse("2006-01-02", y+"-"+md)
	if err != nil {
		return "", err
	}
	_, m, d := date.AddDate(0, 0, 6).Date()

	return fmt.Sprintf("%02d", int(m)) + "-" + fmt.Sprintf("%02d", d), nil
}

func copyDir(srcDir, dstDir string) error {
	entries, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		fileInfo, err := os.Stat(srcPath)
		if err != nil {
			return err
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := createIfNotExists(dstPath); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		default:
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	in, err := os.Open(srcFile)
	if err != nil {
		return err
	}

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return nil
}

func createIfNotExists(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return err
			}
		} else {
			return err
		}
	}

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

func reverse(posts []Post) []Post {
	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}
	return posts
}

func (config *Config) createFileFromTemplate(tmpl *template.Template, dstPath string, data interface{}) error {
	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}

	return tmpl.Execute(f, data)
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
	if err := config.createFileFromTemplate(t, filepath.Join(dir, "index.html"), config); err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	done := make(chan error, 1)
	go func() {
		http.Handle("/", http.FileServer(http.Dir(dir)))
		done <- http.ListenAndServe(":8080", nil)
	}()

	select {
	case <-c:
		return nil
	case err := <-done:
		return err
	}
}

func (config *Config) Publish() error {
	if err := config.Build(config.PublishDir); err != nil {
		return err
	}

	t, err := generateTemplate("index", filepath.Join(config.TemplateDir, "index.html.tmpl"))
	if err != nil {
		return err
	}
	if err := config.createFileFromTemplate(t, filepath.Join(config.PublishDir, "index.html"), config); err != nil {
		return err
	}

	return nil
}
