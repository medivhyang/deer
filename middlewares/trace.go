package middlewares

import (
	"log"
	"net/http"
)

func Trace(callback ...func(r *http.Request)) Middleware {
	var f func(r *http.Request)
	if len(callback) > 0 {
		f = callback[0]
	}
	if f == nil {
		f = func(r *http.Request) {
			log.Printf("%s %s", r.Method, r.URL.Path)
		}
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			f(r)
			h.ServeHTTP(w, r)
		})
	}
}
