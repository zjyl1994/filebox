package service

import (
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func ProcMarkdown(inputFile, baseDir string) (string, error) {
	inputData, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return "", err
	}
	inputData = markdown.NormalizeNewlines(inputData)

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	doc := p.Parse(inputData)

	doc = markdownAddSignOnAst(doc, baseDir)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	html := markdown.Render(doc, renderer)
	return string(html), nil
}

func markdownAddSignOnAst(doc ast.Node, baseDir string) ast.Node {
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if img, ok := node.(*ast.Image); ok && entering {
			img.Destination = procLink(img.Destination, baseDir)
		}

		if link, ok := node.(*ast.Link); ok && entering {
			link.Destination = procLink(link.Destination, baseDir)
		}

		return ast.GoToNext
	})
	return doc
}

func procLink(input []byte, baseDir string) []byte {
	isRelativeUrl := func(link string) bool {
		u, err := url.Parse(link)
		if err != nil {
			return false
		}
		return len(u.Host) == 0 && len(u.Path) > 0
	}
	dst := string(input)
	if isRelativeUrl(dst) {
		dst = strings.TrimPrefix(dst, "./")
		if !strings.HasPrefix(dst, "/") {
			if baseDir == "/" {
				dst = "/" + dst
			} else {
				dst = baseDir + "/" + dst
			}
		}
		dst = dst + "?" + GenSecureLinkStr(dst)
	}
	return []byte(dst)
}
