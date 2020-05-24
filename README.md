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
|- public # path/to/publish
|- content
|  |- mm-dd.md # weekly (default)
|- static
|  |- github-markdown.css
|  |- highlight.min.js
|- template
|  |- index.html.tmpl
|  |- content.md.tmpl
|- config.yaml
```

sample `config.yaml`

```
name: new-project-name
description: This is new project
type: week
```

You can add and delete files in the `static` and `template`, but cannot change directory structure.

### gournal post

You can create new journal post using `gournal post`.

#### Example

`gournal post` at `2020-05-24`

- TypeDaily
  - create `content/2020/05/24.md`
- TypeWeekly
  - create `content/2020/05-18.md` (starting weekday is Monday)
- TypeMonthly
  - create `content/202005.md`

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
|- public # path/to/publish
|- content
|  |- mm-dd.md # weekly
|- static
|  |- github-markdown.css
|  |- highlight.min.js
|- template
|  |- index.html.tmpl
|  |- content.md.tmpl
|- config.yaml
```

you get the build result like below.

```
path/to/publish
|- yyyy
|  |- mm-dd.html
|- static
|  |- github-markdown.css
|  |- highlight.min.js
|- index.html
```

If you set the journal type as monthly, html file is like below.

```
path/to/publish
|- yyyy-mm.html
```

If you set the journal type as daily, html file is like below.

```
path/to/publish
|- yyyy
|  |- mm
|     |- dd.html
```

Default build destination dir is `path/to/project/public`

# LICENSE

MIT
