package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Server struct {
	Port            int
	Watch           bool
	ContentRootPath string
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request URL: %s\n", r.URL.Path)
	p := filepath.Join(s.ContentRootPath, r.URL.Path)
	stat, err := os.Stat(p)
	if err != nil {
		fmt.Print(err)
		return
	}

	if stat.IsDir() {
		entries, err := os.ReadDir(p)

		if err != nil {
			fmt.Print(err)
			return
		}

		var files []fs.FileInfo

		for _, entry := range entries {
			f, err := entry.Info()
			if err != nil {
				fmt.Print(err)
				return
			}

			files = append(files, f)
		}

		writeDirHTML(w, files)
	} else {
		// file
		fmt.Fprintf(w, "%s is file", r.URL.Path)
	}

}

func (s *Server) Start() error {
	http.HandleFunc("/", s.handler)
	fmt.Println("http server started on :" + strconv.Itoa(s.Port))

	if err := http.ListenAndServe(":"+strconv.Itoa(s.Port), nil); err != nil {
		return err
	}

	return nil
}
