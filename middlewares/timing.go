package middlewares

import (
	"github.com/medivhyang/deer"
	"log"
	"net/http"
	"time"
)

func Timing(callback ...func(r *deer.Request, d time.Duration)) deer.Middleware {
	var f func(r *deer.Request, d time.Duration)
	if len(callback) > 0 {
		f = callback[0]
	}
	if f == nil {
		f = func(r *deer.Request, d time.Duration) {
			log.Printf("timing: \"%s %s\" cost %s\n", r.Method(), r.Path(), d)
		}
	}
	return func(h http.Handler) http.Handler {
		return deer.HandlerFunc(func(w *deer.ResponseWriter, r *deer.Request) {
			defer func(start time.Time) {
				d := time.Since(start)
				f(r, d)
			}(time.Now())
			h.ServeHTTP(w.Raw, r.Raw)
		})
	}
}
