package templates

import (
	"html/template"
	"path"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

func Configure(theme string) (*template.Template, error) {
	// This is just a starting point so there's a reasonable policy
	p := bluemonday.UGCPolicy()

	return template.New("").Funcs(template.FuncMap{
		"isodatefmt": func(t time.Time) string {
			return t.Format(time.RFC3339)
		},
		"datefmt": func(fmt string, t time.Time) string {
			return t.Format(fmt)
		},
		"safe": func(text string) template.HTML {
			return template.HTML(text)
		},
		"sanitize": func(text template.HTML) template.HTML {
			return template.HTML(p.Sanitize(string(text)))
		},
		"excerpt": func(max int, text template.HTML) template.HTML {
			return template.HTML(excerpt(string(text), max))
		},
	}).ParseFiles(path.Join(theme, "index.html"))
}
