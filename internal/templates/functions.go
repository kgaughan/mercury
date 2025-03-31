package templates

import (
	"html/template"
	"io/fs"
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
			return template.HTML(text) //nolint:gosec
		},
		"sanitize": func(text template.HTML) template.HTML {
			return template.HTML(p.Sanitize(string(text))) //nolint:gosec
		},
		"excerpt": func(maxlen int, text template.HTML) template.HTML {
			return template.HTML(excerpt(string(text), maxlen)) //nolint:gosec
		},
	}).Funcs(sprig.FuncMap())
}

func Configure(themeFS fs.FS) (*template.Template, error) {
	//nolint:wrapcheck
	return configureFunctions().ParseFS(themeFS, "*.html")
}
