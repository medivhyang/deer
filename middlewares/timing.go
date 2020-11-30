package middlewares

import (
	"github.com/medivhyang/deer"
	"log"
	"time"
)

func Timing(callback ...func(w deer.ResponseWriter, r *deer.Request, d time.Duration)) deer.Middleware {
	var f func(w deer.ResponseWriter, r *deer.Request, d time.Duration)
	if len(callback) > 0 {
		f = callback[0]
	}
	if f == nil {
		f = func(w deer.ResponseWriter, r *deer.Request, d time.Duration) {
			log.Printf("timing: \"%s %s\" cost %s\n", r.Method(), r.Path(), d)
		}
	}
	return func(h deer.HandlerFunc) deer.HandlerFunc {
		return func(w deer.ResponseWriter, r *deer.Request) {
			defer func(start time.Time) {
				d := time.Since(start)
				f(w, r, d)
			}(time.Now())
			h.ServeHTTP(w.Raw(), r.Raw())
		}
	}
}
