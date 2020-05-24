package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type JournalType int

const (
	TypeDaily JournalType = iota
	TypeWeekly
	TypeMonthly
)

type Config struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Type        JournalType `yaml:"type"`
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
		{"", "config.yaml", configTmpl + "type: " + config.typeToString() + "\n"},
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
