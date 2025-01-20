package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
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

var (
	//go:embed web
	web      embed.FS
	tmplFS   fs.FS
	staticFS fs.FS
)

func init() {
	var err error

	tmplFS, err = fs.Sub(web, "web/template")
	if err != nil {
		panic(err)
	}

	staticFS, err = fs.Sub(web, "web/static")
	if err != nil {
		panic(err)
	}
}

type Server struct {
	Port               int
	Watch              bool
	ContentRootPath    string
	ContentRootDirName string
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/"+systemToken) {
		return
	}

	// contentPath is ContentRootPath + r.URL.PathrelativePath)
	contentPath := filepath.Join(s.ContentRootPath, r.URL.Path)
	stat, err := os.Stat(contentPath)
	if err != nil {
		fmt.Print(err)
		return
	}

	if stat.IsDir() {
		entries, err := os.ReadDir(contentPath)
		if err != nil {
			fmt.Print(err)
			return
		}

		d := DirTmplContext{}

		// Add root path
		d.Paths = append(d.Paths, PathInfo{
			Name:     s.ContentRootDirName,
			FullPath: "/",
		})

		// relativePath is the relative path based on the ContentRootPath.
		// dir-root(ContentRootPath)
		// - dir1
		//   - dir2
		// than relativePath is dir1/dir2
		relativePath, _ := strings.CutPrefix(r.URL.Path, "/")

		// if r.URL.Path == "/" than relativePaths == [] so for is not working
		relativePaths := strings.Split(relativePath, "/")

		for i, name := range relativePaths {
			fullpath, err := url.JoinPath("/", strings.Join(relativePaths[:i], "/"), name)
			if err != nil {
				fmt.Println(err)
				return
			}

			d.Paths = append(d.Paths, PathInfo{
				Name:     name,
				FullPath: fullpath,
			})
		}

		for _, entry := range entries {
			f, err := entry.Info()
			if err != nil {
				fmt.Println(err)
				return
			}

			fullpath, err := url.JoinPath("/", r.URL.Path, entry.Name())
			if err != nil {
				fmt.Println(err)
				return
			}

			d.Files = append(d.Files, FileInfo{
				Name:     f.Name(),
				Date:     f.ModTime(),
				IsDir:    f.IsDir(),
				FullPath: fullpath,
			})
		}

		d.Write(w)
	} else {
		// file
		fmt.Fprintf(w, "%s is file", r.URL.Path)
	}

}

func (s *Server) Start() error {
	http.HandleFunc("/", s.handler)
	http.Handle("/"+systemToken+"/", http.StripPrefix("/"+systemToken+"/", http.FileServer(http.FS(staticFS))))

	fmt.Println("http server started on :" + strconv.Itoa(s.Port))

	if err := http.ListenAndServe(":"+strconv.Itoa(s.Port), nil); err != nil {
		return err
	}

	return nil
}
