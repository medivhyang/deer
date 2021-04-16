package middlewares

import (
	"github.com/medivhyang/deer"
)

func MaxAllowed(n int) deer.Middleware {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }
	return func(f deer.HandlerFunc) deer.HandlerFunc {
		acquire()
		defer release()
		return func(w deer.ResponseWriter, r *deer.Request) {
			f.Next(w, r)
		}
	}
}
