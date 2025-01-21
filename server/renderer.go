package server

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/h2non/filetype"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func Render(contentPath, relativePath string) (template.HTML, error) {
	file, err := os.Open(contentPath)
	if err != nil {
		return "", err
	}

	head := make([]byte, 261)
	file.Read(head)

	switch {
	case filetype.IsApplication(head):
		return "", nil
	case filetype.IsArchive(head):
		return "", nil
	case filetype.IsAudio(head):
		return "", nil
	case filetype.IsDocument(head):
		return "", nil
	case filetype.IsFont(head):
		return "", nil
	case filetype.IsImage(head):
		return imgRender(relativePath), nil
	case filetype.IsVideo(head):
		return "", nil
	}

	// recognize as text type
	ext := filepath.Ext(contentPath)
	b, err := os.ReadFile(contentPath)
	if err != nil {
		return "", err
	}

	switch ext {
	case ".markdown", ".md":
		return markdownRender(b)
	}

	return "", nil
}

func imgRender(relativePath string) template.HTML {
	return template.HTML(heredoc.Docf(`
		<img src="/%s/%s" />
	`, staticToken, relativePath))
}

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.CJK,
		meta.New(meta.WithTable()),
	),

	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
	),
)

func markdownRender(source []byte) (template.HTML, error) {
	var buf bytes.Buffer

	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}
	return template.HTML(buf.Bytes()), nil
}
