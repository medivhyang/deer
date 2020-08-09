package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func Recovery(callback ...func(err interface{}, w http.ResponseWriter, r *http.Request)) func(http.Handler) http.Handler {
	var f func(err interface{}, w http.ResponseWriter, r *http.Request)
	if len(callback) > 0 {
		f = callback[0]
	}
	if f == nil {
		f = func(err interface{}, w http.ResponseWriter, r *http.Request) {
			log.Printf("deer: catch panic: %v\n", err)
			log.Println(string(debug.Stack()))
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		}
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					f(err, w, r)
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}
