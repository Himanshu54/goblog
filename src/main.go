package main

import (
	"html/template"
	"log"
	"net/http"

	"app/posts"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /home", NewHomePage("header", "footer", "content").Handler())
	mux.HandleFunc("GET /posts/", posts.NewPostPage().Handler())
	mux.HandleFunc("GET /posts/{slug}", posts.NewPostPage().SlugHandler())

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
