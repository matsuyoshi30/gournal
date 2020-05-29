package main

var indexTmpl = `<!DOCTYPE html>
<html>
  <head>
    <title>{{ .Name }}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/4.0.0/github-markdown.min.css">
    <link rel="stylesheet" href="./static/styles.css">
  </head>
  <body>
    <article class="markdown-body">
      <h1>{{ .Name }}</h1>
      {{ range $i, $p := .Posts }}
      {{ if or ( eq $i 0 ) ( $p.IsLastWeek ) }}<h2>{{ $p.PostYear }}</h2>{{ end }}
      <ul>
        <li><a href="{{ $p.Link }}">Week {{ $p.WeekNum }}</a> <small>{{ $p.FromDate }} - {{ $p.ToDate }}</small></li>
      </ul>
      {{ end }}
    </article>
  </body>
</html>
`

var cssTmpl = `.markdown-body {
  box-sizing: border-box;
  min-width: 200px;
  max-width: 980px;
  margin: 0 auto;
  padding: 45px;
}
@media (max-width: 767px) {
  .markdown-body {
    padding: 15px;
  }
}
@media (prefers-color-scheme: dark) {
  body {
    background-color: #1b262c;
  }
  .markdown-body img {
    background-color: #1b262c;
  }
  .markdown-body {
    color: #e0e0e0;
  }
  .markdown-body blockquote {
    color: #95989a;
    border-left: .25em solid #6d6d6d;
  }
  .markdown-body a {
    color: #6bb4ff;
  }
}
`

var contentTmpl = `## Title

### Contents

`

var postTmpl = `<!DOCTYPE html>
<html>
  <head>
    <title>{{ .Title }}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/4.0.0/github-markdown.min.css">
    <link rel="stylesheet" href="{{ .CSSPath }}">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.18.1/styles/default.min.css">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.18.1/highlight.min.js"></script>
    <script>hljs.initHighlightingOnLoad();</script>
  </head>
  <body>
    <article class="markdown-body">
      <h1>Week {{ .WeekNum }} - {{ .PostYear }}</h1>
      {{ if ne .FromDate "" }}
        The weekly report from {{ .FromDate }} to {{ .ToDate }}
      {{ end }}
      {{ .Body }}
    </article>
  </body>
</html>
`

var configTmpl = `name: your project name
description: your project description
`
