package templates

import (
	"html/template"
	"path"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/microcosm-cc/bluemonday"
)

func configureFunctions() *template.Template {
	// This is just a starting point so there's a reasonable policy
	p := bluemonday.UGCPolicy()

	return template.New("").Funcs(template.FuncMap{
		"isodate": func(t time.Time) string {
			return t.Format(time.RFC3339)
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
	}).Funcs(sprig.FuncMap())
}

func Configure(theme string) (*template.Template, error) {
	return configureFunctions().ParseFiles(path.Join(theme, "index.html"))
}
