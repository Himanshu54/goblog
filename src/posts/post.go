package posts

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Post struct {
	Title   string
	Tags    []string
	Content template.HTML
}

func (p *Post) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		fmt.Fprint(w, "post: ", slug)
	}
}

type PostPage struct {
	Posts []PostData
}

type PostData struct {
	Title    string    `yaml:"title"`
	Tags     []string  `yaml:"tags"`
	Date     time.Time `yaml:"date"`
	PostFile string
}

type PostContent struct {
	Data  template.HTML
	Title string
}

func NewPostPage() *PostPage {
	return &PostPage{}
}

func (pp *PostPage) Handler() http.HandlerFunc {
	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"date": func(x time.Time) string { return x.Format(time.RFC822) },
	}
	return func(w http.ResponseWriter, r *http.Request) {
		postPage := PostPage{}
		entries, err := os.ReadDir("posts/content/")
		// fmt.Print(entries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".md") {
				f, err := os.Open("posts/content/" + e.Name())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer f.Close()

				matter := PostData{}
				_, err = frontmatter.MustParse(f, &matter)
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				postPage.Posts = append(postPage.Posts, PostData{Title: matter.Title, Tags: matter.Tags, Date: matter.Date, PostFile: e.Name()})
			}
		}
		if len(postPage.Posts) == 0 {
			http.Error(w, "no post to show", http.StatusNotFound)
			return
		}
		postTemp, err := template.New(path.Base("posts/postpage.html")).Funcs(funcMap).ParseFiles("posts/postpage.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = postTemp.Execute(w, postPage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (pp *PostPage) SlugHandler() http.HandlerFunc {
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
	sl := FileReader{}
	tmpl := template.Must(template.ParseFiles("posts/post.html"))
	return func(w http.ResponseWriter, r *http.Request) {

		slug := r.PathValue("slug")
		postMarkdown, err := sl.Read(slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		matter := PostData{}
		postContent := PostContent{}
		rest, err := frontmatter.MustParse(strings.NewReader(postMarkdown), &matter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		var buf bytes.Buffer
		if err := md.Convert(rest, &buf); err != nil {
			panic(err)
		}
		postContent.Data = template.HTML(buf.String())
		postContent.Title = matter.Title
		tmpl.Execute(w, postContent)
	}
}

type SlugReader interface {
	Read(slug string) (string, error)
}

type FileReader struct{}

func (fsr FileReader) Read(slug string) (string, error) {
	f, err := os.Open("posts/content/" + slug)
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return "", nil
	}
	return string(b), nil
}
