# gournal - Go + Journal

## Usage

### gournal new

```
$ gournal new [-month|-week(default)|-day] <projectname>
```

If you command `gournal new week newproject`, make the project directory like below.

```
.newproject
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
└── config.json
```

sample `config.json`

```
{
  "name": "newproject",
  "version": "1.0.0",
  "description": "weekly journal",
  "type": "weekly",
  "title": "weekly journal",
  "domain": "",
  "repository": {
    "type": "git",
    "url": "git+https://github.com/matsuyoshi30/gournal.git"
  },
  "author": "matsuyoshi30",
  "license": "MIT",
  "bugs": {
    "url": "https://github.com/matsuyoshi30/gournal/issues"
  },
  "homepage": "https://github.com/matsuyoshi30/gournal#readme"
}
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
.project root
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
.build destination dir
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

Default build destination dir is `project/build`

# LICENSE

MIT