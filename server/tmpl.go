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
	Paths       []PathInfo
	Files       []FileInfo
	SystemToken string
}

type FileTmplContext struct {
	Paths       []PathInfo
	SystemToken string
	HTML        template.HTML
}

var (
	dirTmpl  *template.Template
	fileTmpl *template.Template
)

func init() {
	b, err := fs.ReadFile(tmplFS, "dir.html")
	if err != nil {
		panic(err)
	}

	dirTmpl = template.Must(template.New("dir-tmpl").Parse(string(b)))

	b, err = fs.ReadFile(tmplFS, "file.html")
	if err != nil {
		panic(err)
	}

	fileTmpl = template.Must(template.New("file-tmpl").Parse(string(b)))
}

func (d *DirTmplContext) Write(w http.ResponseWriter) {
	if err := dirTmpl.Execute(w, d); err != nil {
		fmt.Println(err)
	}
}

func (f *FileTmplContext) Write(w http.ResponseWriter) {
	if err := fileTmpl.Execute(w, f); err != nil {
		fmt.Println(err)
	}
}
