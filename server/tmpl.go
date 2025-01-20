package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"time"
)

type PathInfo struct {
	Name     string
	FullPath string
}

type FileInfo struct {
	Name        string
	FullPath    string
	Description string
	IsDir       bool
	Date        time.Time
}

type DirTmplContext struct {
	Paths []PathInfo
	Files []FileInfo
}

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

func (d *DirTmplContext) Write(w http.ResponseWriter) {
	if err := dirTmpl.Execute(w, d); err != nil {
		fmt.Println(err)
	}
}
