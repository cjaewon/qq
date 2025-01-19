package server

import (
	"fmt"
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

		for _, entry := range entries {
			fmt.Fprint(w, entry.Name()+"<br/>")
		}
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
