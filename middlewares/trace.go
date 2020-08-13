package middlewares

import (
	"github.com/medivhyang/deer"
	"log"
	"net/http"
)

func Trace(callback ...func(r *deer.Request)) deer.Middleware {
	var f func(r *deer.Request)
	if len(callback) > 0 {
		f = callback[0]
	}
	if f == nil {
		f = func(r *deer.Request) {
			log.Printf("%s %s", r.Method(), r.Path())
		}
	}
	return func(h http.Handler) http.Handler {
		return deer.HandlerFunc(func(w *deer.ResponseWriter, r *deer.Request) {
			f(r)
			h.ServeHTTP(w.Raw, r.Raw)
		})
	}
}
