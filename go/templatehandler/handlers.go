package templatehandler

import (
	"fmt"
	"net/http"
)

// dynamicHandler serves responses based on the http request.
type dynamicHandler struct {
	t *TemplateHandler
	f func(w http.ResponseWriter, r *http.Request) map[string]interface{}
}

func (d *dynamicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := d.t.render(w, r, d.f(w, r))
	if err != nil {
		return
	}
	w.Write(b)
}

func (t *TemplateHandler) Dynamic(f func(w http.ResponseWriter, r *http.Request) map[string]interface{}) *dynamicHandler {
	return &dynamicHandler{
		t: t,
		f: f,
	}
}

// staticHandler serves responses based on a provided map of input (or nil), and caches the response.
type staticHandler struct {
	t *TemplateHandler
	m map[string]interface{}
	c []byte
}

func (s *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.c == nil {
		b, err := s.t.render(w, r, s.m)
		if err != nil {
			panic(fmt.Errorf("could not render static template: %v", err))
		}
		s.c = b
	}

	w.Write(s.c)
}

func (t *TemplateHandler) Static(m map[string]interface{}) *staticHandler {
	return &staticHandler{
		t: t,
		m: m,
	}
}
