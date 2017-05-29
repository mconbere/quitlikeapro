// Package templatehandler looks complicated, but in fact is a relatively simple wrapper around Go's html/template.
//
// The concept is that you should be able to create a simple webpage that consists of a base "wrapper" template that
// contains the header and footer html. You then can provide the following html templates for each custom page:
//
// - "content": This is the main body of your page.
// - "css": This is any additional CSS you want to add. It's optional, and added to the bottom of the existing CSS.
// - "js": This is any additional Javascript you want to add. It's optional, and added to the bottom of the existing Javascript.
// - "input": This is a JSON blob. Here you can add custom elements to the html template's pipeline.
//
// Here is a simple example for rendering an index.html with a base.html:
//
//     base.html:
//     {{ define "base" }}
//     <!doctype html>
//     <html>
//       <head>
//         <title>{{ .title }}</title>
//         {{ template "css" . }}
//       </head>
//       <body>
//         <h1>{{ .title }}</h1>
//         {{ template "content" . }}
//         {{ template "js" . }}
//       </body>
//     </html>
//     {{ end }}
//
//     index.html:
//     {{ define "input" }}
//     {
//       "title": "Index"
//     }
//     {{ end }}
//     {{ define "content" }}
//     Some content.
//     {{ end }}
//
//     main.go:
//     b, _ := NewBase("base.html", nil)
//     h, _ := New(b, "/", "index.html", nil)
//     http.Handle("/", h)
package templatehandler

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"
)

type Base struct {
	Template *template.Template
	Input    map[string]interface{}
}

func NewBase(tmpl string, input map[string]interface{}) (*Base, error) {
	t := template.New("")
	t, err := t.ParseFiles(tmpl)
	if err != nil {
		return nil, err
	}
	return &Base{
		Template: t,
		Input:    input,
	}, nil
}

type TemplateHandler struct {
	Template *template.Template
	Input    map[string]interface{}
}

func New(base *Base, tmpl string) (*TemplateHandler, error) {
	t, err := base.Template.Clone()
	if err != nil {
		return nil, err
	}

	t.Funcs(template.FuncMap{
		"markdown": Markdown(t),
	})

	t, err = t.ParseFiles(tmpl)
	if err != nil {
		return nil, err
	}

	if t.Lookup("js") == nil {
		if _, err := t.Parse("{{ define \"js\" }}{{ end }}"); err != nil {
			return nil, err
		}
	}
	if t.Lookup("css") == nil {
		if _, err := t.Parse("{{ define \"css\" }}{{ end }}"); err != nil {
			return nil, err
		}
	}

	input := base.Input
	if t.Lookup("input") != nil {
		js, err := jsonFromTmpl(t, "input")
		if err != nil {
			return nil, err
		}
		input = mergeMap(base.Input, js)
	}

	return &TemplateHandler{
		Template: t,
		Input:    input,
	}, nil
}

func Must(t *TemplateHandler, err error) *TemplateHandler {
	if err != nil {
		panic(err)
	}
	return t
}

func cleanPath(in string) string {
	out := strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(in, "index.html"), "/"))
	if out == "" {
		out = "/"
	}
	return out
}

func (t *TemplateHandler) render(w http.ResponseWriter, r *http.Request, input map[string]interface{}) ([]byte, error) {
	input = mergeMap(t.Input, input)

	var b bytes.Buffer
	err := t.Template.ExecuteTemplate(&b, "base", input)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
