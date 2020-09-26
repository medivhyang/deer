package middlewares

import (
	"github.com/medivhyang/deer"
	"log"
)

func Trace(callback ...func(w *deer.ResponseWriter, r *deer.Request)) deer.Middleware {
	var f func(w *deer.ResponseWriter, r *deer.Request)
	if len(callback) > 0 {
		f = callback[0]
	}
	if f == nil {
		f = func(w *deer.ResponseWriter, r *deer.Request) {
			log.Printf("%s %s", r.Method(), r.Path())
		}
	}
	return func(h deer.HandlerFunc) deer.HandlerFunc {
		return func(w *deer.ResponseWriter, r *deer.Request) {
			f(w, r)
			h.ServeHTTP(w.Raw, r.Raw)
		}
	}
}
