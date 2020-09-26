package middlewares

import (
	"github.com/medivhyang/deer"
	"log"
	"time"
)

func Timing(callback ...func(d time.Duration, w *deer.ResponseWriter, r *deer.Request)) deer.Middleware {
	var f func(d time.Duration, w *deer.ResponseWriter, r *deer.Request)
	if len(callback) > 0 {
		f = callback[0]
	}
	if f == nil {
		f = func(d time.Duration, w *deer.ResponseWriter, r *deer.Request) {
			log.Printf("timing: \"%s %s\" cost %s\n", r.Method, r.Path(), d)
		}
	}
	return func(h deer.HandlerFunc) deer.HandlerFunc {
		return func(w *deer.ResponseWriter, r *deer.Request) {
			defer func(start time.Time) {
				d := time.Since(start)
				f(d, w, r)
			}(time.Now())
			h.ServeHTTP(w.Raw, r.Raw)
		}
	}
}
