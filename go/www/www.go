package www

import (
	"net/http"

	"html/template"

	"github.com/mconbere/quitlikeapro/go/templatehandler"
)

type Quittable struct {
	Title template.HTML
	Steps []template.HTML
}

var quittables = []Quittable{{
	Title: "Emacs",
	Steps: []template.HTML{
		"Hold down <code>CTRL</code>",
		"Press <code>x</code>",
		"Press <code>c</code>",
	},
}, {
	Title: "Vim",
	Steps: []template.HTML{
		"Type <code>:q</code>",
		"Press <code>enter</code>",
	},
}, {
	Title: "Python Interpreter",
	Steps: []template.HTML{
		"Type <code>CTRL</code>-<code>d</code>",
	},
}, {
	Title: "Every other command line tool",
	Steps: []template.HTML{
		"Type <code>CTRL</code>-<code>c</code>",
	},
}}

func New() http.Handler {
	mux := http.NewServeMux()

	base, err := templatehandler.NewBase("templates/base.html", map[string]interface{}{
		"Quittables": quittables,
	})
	if err != nil {
		panic(err)
	}

	mux.Handle("/about", templatehandler.Must(templatehandler.New(base, "templates/about/index.html")).Static(nil))
	mux.Handle("/", templatehandler.Must(templatehandler.New(base, "templates/index.html")).Static(nil))

	return mux
}
