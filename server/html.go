package server

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

var (
	//go:embed web
	web embed.FS

	dirTmpl *template.Template
)

func init() {
	b, err := web.ReadFile("./html/dir.html")
	if err != nil {
		panic(err)
	}

	dirTmpl = template.Must(template.New("dir-tmpl").Parse(string(b)))
}

func writeDirHTML(w http.ResponseWriter, files []fs.FileInfo) {
	if err := dirTmpl.Execute(w, files); err != nil {
		fmt.Println(err)
	}
}
