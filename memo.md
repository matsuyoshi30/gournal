`gournal new`

1. check working directory
1. make directory structure
1. download `github-markdown.css` and `highlight.min.js` from original
1. make template file
1. make content file
1. make config file

`gournal test`

1. run local html server
1. check working directory structure and files
1. make directory structure into destination directory
1. read template
1. read contents and make html file from template and content
1. make index.html
1. in the end, delete all build result

continues to logging when user is editting contents or static files.

正常にビルドできていることを返り値で表現して CI で gournal test -> if test pass, gournal pub みたいなことがしたい

`gournal pub`

1. check working directory structure and files
1. make directory structure into destination directory
1. read template
1. read contents and make html file from template and content
1. make index.html
