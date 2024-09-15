package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Site struct {
	indexPage string
}

func (s *Site) handler(content template.HTML) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Site) IndexHandler() http.HandlerFunc {
	return s.handler(template.HTML(""))
}

func main() {
	site := Site{
		indexPage: "index.html",
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", site.IndexHandler())
	mux.HandleFunc("GET /{page}/{slug}", site.PageHandler())

	dir := http.Dir("./static")
	fs := http.FileServer(dir)
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	wrappedMux := NewLoggingMiddleWare(mux)
	err := http.ListenAndServe(":8080", wrappedMux)
	if err != nil {
		log.Fatal(err)
	}
}

type HomePage struct {
	Header  template.HTML
	Footer  template.HTML
	Content template.HTML
}

func NewHomePage(header, footer, content string) *HomePage {
	return &HomePage{
		Header:  template.HTML(header),
		Footer:  template.HTML(footer),
		Content: template.HTML(content),
	}
}

func (hp *HomePage) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		if err := tmpl.Execute(w, *hp); err != nil {
			panic(err)
		}
	}
}

func NewScribblePageHandler(w http.ResponseWriter, r *http.Request) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai")),
		),
	)
	scribble, err := os.Open("scribble.md")
	if err != nil {
		panic(err)
	}
	scribbleText, err := io.ReadAll(scribble)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	md.Convert(scribbleText, &buf)
	tmpl := template.Must(template.ParseFiles("index.html"))
	if err := tmpl.Execute(w, NewHomePage("header", "footer", buf.String())); err != nil {
		panic(err)
	}
}
