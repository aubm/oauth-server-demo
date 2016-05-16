package api

import "net/http"

type Adapter interface {
	Adapt(next http.Handler) http.Handler
}

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter.Adapt(h)
	}
	return h
}
