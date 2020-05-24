package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

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
	Type        JournalType
}

var config Config

var indexTmpl = `<!DOCTYPE html>
<html>
  <head>
    <title></title>
  </head>
  <body>
    <h1></h1>
  </body>
</html>
`

var contentTmpl = `# Title

## Contents

`

var configTmpl = `name: <your project name>
description: <your project description>
`

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

		if err := ioutil.WriteFile(file, []byte(f.contents), 0755); err != nil {
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
	//   TypeDaily   -> `content/2020/05/24.md`
	//   TypeWeekly  -> `content/2020/05-18.md` check starting weekday
	//   TypeMonthly -> `content/202005.md`
	dir := filepath.Join(config.Wd, "content")
	if config.Type == TypeDaily {
		dir = filepath.Join(dir, strconv.Itoa(year), monthNameToNum(month))
	} else if config.Type == TypeWeekly {
		dir = filepath.Join(dir, strconv.Itoa(year))
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// read content template
	b, err := ioutil.ReadFile(filepath.Join(config.Wd, "template", "content.md.tmpl"))
	if err != nil {
		return err
	}

	// create post file
	var filename string
	if config.Type == TypeDaily {
		filename = strconv.Itoa(day) + ".md"
	} else if config.Type == TypeWeekly {
		filename = checkWeekday(time.Now()) + ".md"
	} else if config.Type == TypeMonthly {
		filename = strconv.Itoa(year) + month.String() + ".md"
	}
	if err := ioutil.WriteFile(filepath.Join(dir, filename), b, 0755); err != nil {
		return err
	}

	return nil
}

func monthNameToNum(month time.Month) string {
	var monthNum string
	switch month {
	case time.January:
		monthNum = "01"
	case time.February:
		monthNum = "02"
	case time.March:
		monthNum = "03"
	case time.April:
		monthNum = "04"
	case time.May:
		monthNum = "05"
	case time.June:
		monthNum = "06"
	case time.July:
		monthNum = "07"
	case time.August:
		monthNum = "08"
	case time.September:
		monthNum = "09"
	case time.October:
		monthNum = "10"
	case time.November:
		monthNum = "11"
	case time.December:
		monthNum = "12"
	}

	return monthNum
}

func checkWeekday(date time.Time) string {
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
	}
	_, month, day := date.Date()

	return monthNameToNum(month) + "-" + strconv.Itoa(day)
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
	config.Type = config.stringToType()

	return nil
}
