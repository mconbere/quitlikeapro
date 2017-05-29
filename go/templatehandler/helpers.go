package templatehandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/russross/blackfriday"
)

func Markdown(t *template.Template) func(string, interface{}) (template.HTML, error) {
	return func(name string, in interface{}) (template.HTML, error) {
		var b bytes.Buffer
		if err := t.ExecuteTemplate(&b, name, in); err != nil {
			return "", err
		}
		out := blackfriday.MarkdownCommon(b.Bytes())
		return template.HTML(out), nil
	}
}

func mergeMap(base, in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range base {
		out[k] = v
	}
	for k, v := range in {
		out[k] = v
	}
	return out
}

func jsonFromTmpl(t *template.Template, name string) (map[string]interface{}, error) {
	var base bytes.Buffer
	t.ExecuteTemplate(&base, name, nil)
	input := make(map[string]interface{})
	if err := json.Unmarshal(base.Bytes(), &input); err != nil {
		return nil, fmt.Errorf("json unmarshaling of template %q failed: %v", name, err)
	}
	return input, nil
}
