package middlewares

import (
	"log"
	"net/http"
	"time"
)

func Timing(callback ...func(r *http.Request, d time.Duration)) Middleware {
	var f func(r *http.Request, d time.Duration)
	if len(callback) > 0 {
		f = callback[0]
	}
	if f == nil {
		f = func(r *http.Request, d time.Duration) {
			log.Printf("timing: \"%s %s\" cost %s\n", r.Method, r.URL.Path, d)
		}
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func(start time.Time) {
				d := time.Since(start)
				f(r, d)
			}(time.Now())
			h.ServeHTTP(w, r)
		})
	}
}
