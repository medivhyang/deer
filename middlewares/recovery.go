package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/medivhyang/deer"
)

func Recovery(callback ...func(err interface{}, w *deer.ResponseWriter, r *deer.Request)) deer.Middleware {
	var f func(err interface{}, w *deer.ResponseWriter, r *deer.Request)
	if len(callback) > 0 {
		f = callback[0]
	}
	if f == nil {
		f = func(err interface{}, w *deer.ResponseWriter, r *deer.Request) {
			log.Printf("deer: catch panic: %v\n", err)
			log.Println(string(debug.Stack()))
			w.Text(http.StatusInternalServerError, fmt.Sprint(err))
		}
	}
	return func(h http.Handler) http.Handler {
		return deer.HandlerFunc(func(w *deer.ResponseWriter, r *deer.Request) {
			defer func() {
				if err := recover(); err != nil {
					f(err, w, r)
				}
			}()
			h.ServeHTTP(w.Raw, r.Raw)
		})
	}
}
