package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// systemToken is a token used for WebSocket, CSS, and JS to prevent unintentional access.
	// localhost:1234/styles/common.css might conflict with the user's files,
	// so systemToken injects a token, resulting in URLs like localhost:1234/[token]/styles/common.css
	systemToken = "JjgVf8b1FVbEuyD2aqnewz3I6z-i8bePpdRwfCF6wi1Puz3vSYqsfNlH9XnHX7tXBomLwpISWT0"
)

type Server struct {
	Port            int
	Watch           bool
	ContentRootPath string
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/"+systemToken) {
		return
	}

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
	http.Handle("/"+systemToken+"/", http.StripPrefix("/"+systemToken+"/", http.FileServer(http.FS(web))))

	fmt.Println("http server started on :" + strconv.Itoa(s.Port))

	if err := http.ListenAndServe(":"+strconv.Itoa(s.Port), nil); err != nil {
		return err
	}

	return nil
}
