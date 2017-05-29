package main

import (
	"net/http"

	"github.com/mconbere/quitlikeapro/go/www"
)

func init() {
	root := www.New()
	http.Handle("/", root)
}
