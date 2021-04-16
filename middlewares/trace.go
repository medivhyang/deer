package middlewares

import (
	"github.com/medivhyang/deer"
)

func Trace(callback ...func(w deer.ResponseWriter, r *deer.Request)) deer.Middleware {
	var f func(w deer.ResponseWriter, r *deer.Request)
	if len(callback) > 0 {
		f = callback[0]
	} else {
		f = func(w deer.ResponseWriter, r *deer.Request) {
			logf("%s %s", r.Method(), r.Path())
		}
	}
	return func(h deer.HandlerFunc) deer.HandlerFunc {
		return func(w deer.ResponseWriter, r *deer.Request) {
			f(w, r)
			h.Next(w, r)
		}
	}
}
