package server

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/h2non/filetype"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
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
		return notSupportRender(relativePath), nil
	case filetype.IsArchive(head):
		return notSupportRender(relativePath), nil
	case filetype.IsAudio(head):
		return audioRender(relativePath), nil
	case filetype.IsDocument(head):
		return notSupportRender(relativePath), nil
	case filetype.IsFont(head):
		return notSupportRender(relativePath), nil
	case filetype.IsImage(head):
		return imgRender(relativePath), nil
	case filetype.IsVideo(head):
		return videoRender(relativePath), nil
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

	lex := lexers.Match(filepath.Base(relativePath))
	if lex != nil {
		return codeRender(b, lex)
		// lex pointer check
	}

	return notSupportRender(relativePath), nil
}

func imgRender(relativePath string) template.HTML {
	// relativePath string is r.URL.Path

	return template.HTML(heredoc.Docf(`
		<div class="img-filetype-container">
			<img src="/%s/%s" />
		</div>
	`, staticToken, relativePath))
}

func audioRender(relativePath string) template.HTML {
	return template.HTML(heredoc.Docf(`
		<div class="audio-filetype-container">
			<figure>
				<audio controls src="/%s/%s"></audio>
			</figure>
		</div>
	`, staticToken, relativePath))
}

func videoRender(relativePath string) template.HTML {
	return template.HTML(heredoc.Docf(`
		<div class="video-filetype-container">
			<video controls>
				<source src="/%s/%s" />

				<p>브라우저가 지원하지 않는 비디오 형식입니다. 다른 프로그램을 통해 재생해야합니다.</p>
			</video>
		</div>
	`, staticToken, relativePath))
}

// todo: 다국적
// todo: video 용량, streaming check

var md = goldmark.New(
	// todo: create it as extension
	goldmark.WithExtensions(
		extension.GFM,
		extension.CJK,
		meta.New(meta.WithTable()),
		// todo: 코드 하이라이팅
		// todo: 테이블에 class injection해서 css 적용 수정하기
		// todo: 이미지 url에 특정 staticToken 추가하기
		highlighting.NewHighlighting(
			highlighting.WithStyle("emacs"),
			highlighting.WithFormatOptions(
				chromahtml.WithLineNumbers(true),
				chromahtml.WithClasses(true),
			),
		),
	),

	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)

func markdownRender(source []byte) (template.HTML, error) {
	var buf bytes.Buffer

	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}
	return template.HTML(heredoc.Docf(`
		<div class="markdown-filetype-container">
			%s
		</div>
	`, buf.Bytes())), nil
}

func codeRender(source []byte, lex chroma.Lexer) (template.HTML, error) {
	var buf bytes.Buffer

	if err := quick.Highlight(&buf, string(source), lex.Config().Name, "html", "emacs"); err != nil {
		return "", err
	}

	// todo: background color
	return template.HTML(heredoc.Docf(`
		<div class="code-filetype-container">
			%s
		</div>
	`, buf.String())), nil
}

func notSupportRender(relativePath string) template.HTML {
	base := filepath.Base(relativePath)

	return template.HTML(heredoc.Docf(`
		<div class="notsupport-filetype-container">
			"%s" 는 지원하지 않는 파일 확장자 입니다. <br/>
			전용 프로그램을 사용해주세요.
		</div>
	`, base))
}
