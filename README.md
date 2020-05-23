# gournal - Go + Journal

## Usage

### gournal new

```
$ mkdir path/to/project
$ gournal new [--month|--week(default)|--day] path/to/project
```

If you command `gournal new --week newproject`, make the project directory like below.

```
path/to/project
├─public # path/to/publish
├── content
│   └── yyyy
│       └── mm-dd.md
├── static
│   ├── css
│   │   └── github-markdown.css
│   ├── js
│   │   └── highlight.min.js
│   └── img
├── template
│   ├── index.html.tmpl
│   └── content.html.tmpl
└── confyaml
```

sample `config.yaml`

```
name: new-project-name
description: This is new project
type: week
author: matsuyoshi30
```

You can add and delete files in the `static` and `template`, but cannot change directory structure.

### gournal test

You can check your contents and static files by `gournal test`.

It runs local server, and build html pages.

### gournal pub

```
$ gournal pub
```

Build html page contains static files from markdown and template file.

If you make the project like below and command `gournal pub`,

```
path/to/project
├─public # path/to/publish
├── content
│   ├── yyyy
│   │   └── yy-mm.md  # Monthly
│   ├── yyyy
│   │   └── mm-dd.md  # Weekly
│   └── yyyy
│       └── mm
│           └── dd.md # Daily
├── static  # if necessarily
│   ├── css
│   │   └── stylesheet.css
│   ├── js
│   │   └── index.js
│   ├── img
│   │   └── logo.png
│   ├── CNAME
│   └── favicon.ico
├── template
│   ├── index.html.tmpl
│   └── content.html.tmpl
└── config.json
```

you get the build result like below.

```
path/to/publish
├── yyyy
│   └── yy-mm.md  # Monthly
├── yyyy
│   └── mm-dd.md  # Weekly
├── yyyy
│   └── mm
│       └── dd.md # Daily
├── css
│   └── stylesheet.css
├── js
│   └── index.js
├── img
│   └── logo.png
├── CNAME
├── favicon.ico
└── index.html
```

Default build destination dir is `path/to/project/public`

# LICENSE

MIT
