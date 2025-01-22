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

	// staticToken is a token used for serving static file (video, images, pdf)
	// not for serving web/static, serving for user's file
	staticToken = "MFMEVmXenkenmEkTMemoMywoekQpeLmhretiowktoQETR-qetMlqemltMqjJyeoKypElqQykpq3"
)

var (
	//go:embed web
	web      embed.FS
	tmplFS   fs.FS
	assetsFS fs.FS
)

func init() {
	var err error

	tmplFS, err = fs.Sub(web, "web/template")
	if err != nil {
		panic(err)
	}

	// todo: static->assets change naming
	assetsFS, err = fs.Sub(web, "web/assets")
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
	// todo: study golang http
	if strings.HasPrefix(r.URL.Path, "/"+systemToken) || strings.HasPrefix(r.URL.Path, "/"+staticToken) {
		return
	}

	// relativePath is the relative path based on the ContentRootPath.
	// dir-root(ContentRootPath)
	// - dir1
	//   - dir2
	// than relativePath is dir1/dir2
	relativePath, _ := strings.CutPrefix(r.URL.Path, "/")

	// if r.URL.Path == "/" than relativePaths == [] so for is not working
	relativePaths := strings.Split(relativePath, "/")

	// contentPath is ContentRootPath + r.URL.Path(relativePath)
	contentPath := filepath.Join(s.ContentRootPath, relativePath)
	stat, err := os.Stat(contentPath)
	if err != nil {
		fmt.Print(err)
		return
	}

	if stat.IsDir() {
		entries, err := os.ReadDir(contentPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		d := DirTmplContext{
			SystemToken: systemToken,
		}

		// Add root path
		d.Paths = append(d.Paths, PathInfo{
			Name:     s.ContentRootDirName,
			FullPath: "/",
		})

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
		f := FileTmplContext{
			SystemToken: systemToken,
		}

		// Add root path
		f.Paths = append(f.Paths, PathInfo{
			Name:     s.ContentRootDirName,
			FullPath: "/",
		})

		for i, name := range relativePaths {
			fullpath, err := url.JoinPath("/", strings.Join(relativePaths[:i], "/"), name)
			if err != nil {
				fmt.Println(err)
				return
			}

			f.Paths = append(f.Paths, PathInfo{
				Name:     name,
				FullPath: fullpath,
			})
		}

		f.HTML, err = Render(contentPath, relativePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		f.Write(w)
	}

}

func (s *Server) Start() error {
	http.HandleFunc("/", s.handler)

	http.Handle("/"+systemToken+"/", http.StripPrefix("/"+systemToken+"/", http.FileServer(http.FS(assetsFS))))
	http.Handle("/"+staticToken+"/", http.StripPrefix("/"+staticToken+"/", http.FileServer(http.Dir(s.ContentRootPath))))

	fmt.Println("http server started on :" + strconv.Itoa(s.Port))

	if err := http.ListenAndServe(":"+strconv.Itoa(s.Port), nil); err != nil {
		return err
	}

	return nil
}
