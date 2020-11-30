package middlewares

import (
	"fmt"
	"github.com/medivhyang/deer"
	"log"
	"net/http"
	"runtime/debug"
)

func Recovery(callback ...func(w deer.ResponseWriter, r *deer.Request, err interface{})) deer.Middleware {
	var f func(w deer.ResponseWriter, r *deer.Request,err interface{})
	if len(callback) > 0 {
		f = callback[0]
	}
	if f == nil {
		f = func(w deer.ResponseWriter, r *deer.Request, err interface{}) {
			log.Printf("deer: catch panic: %v\n", err)
			log.Println(string(debug.Stack()))
			http.Error(w.Raw(), fmt.Sprint(err), http.StatusInternalServerError)
		}
	}
	return func(h deer.HandlerFunc) deer.HandlerFunc {
		return func(w deer.ResponseWriter, r *deer.Request) {
			defer func() {
				if err := recover(); err != nil {
					f(w, r, err)
				}
			}()
			h.ServeHTTP(w.Raw(), r.Raw())
		}
	}
}
