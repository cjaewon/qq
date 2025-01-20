package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

var (
	dirTmpl *template.Template
)

func init() {
	b, err := fs.ReadFile(tmplFS, "dir.html")
	if err != nil {
		panic(err)
	}

	dirTmpl = template.Must(template.New("dir-tmpl").Parse(string(b)))
}

type PathInfo struct {
	Name     string
	FullPath string
}

type DirTmplContext struct {
	Paths []PathInfo
	Files []fs.FileInfo
}

func (d *DirTmplContext) Write(w http.ResponseWriter) {
	if err := dirTmpl.Execute(w, d); err != nil {
		fmt.Println(err)
	}
}
