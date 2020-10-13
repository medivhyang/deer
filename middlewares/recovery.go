package middlewares

import (
	"fmt"
	"github.com/medivhyang/deer"
	"log"
	"net/http"
	"runtime/debug"
)

func Recovery(callback ...func(err interface{}, w deer.ResponseWriter, r *deer.Request)) deer.Middleware {
	var f func(err interface{}, w deer.ResponseWriter, r *deer.Request)
	if len(callback) > 0 {
		f = callback[0]
	}
	if f == nil {
		f = func(err interface{}, w deer.ResponseWriter, r *deer.Request) {
			log.Printf("deer: catch panic: %v\n", err)
			log.Println(string(debug.Stack()))
			http.Error(w.Raw(), fmt.Sprint(err), http.StatusInternalServerError)
		}
	}
	return func(h deer.HandlerFunc) deer.HandlerFunc {
		return func(w deer.ResponseWriter, r *deer.Request) {
			defer func() {
				if err := recover(); err != nil {
					f(err, w, r)
				}
			}()
			h.ServeHTTP(w.Raw(), r.Raw())
		}
	}
}
