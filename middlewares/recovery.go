package middlewares

import (
	"fmt"
	"github.com/medivhyang/deer"
	"net/http"
	"runtime/debug"
)

func Recovery(callback ...func(w deer.ResponseWriter, r *deer.Request, err interface{})) deer.Middleware {
	var f func(w deer.ResponseWriter, r *deer.Request, err interface{})
	if len(callback) > 0 {
		f = callback[0]
	} else {
		f = func(w deer.ResponseWriter, r *deer.Request, err interface{}) {
			logf("deer: catch panic: %+v\n%s", err, string(debug.Stack()))
			w.Text(http.StatusInternalServerError, fmt.Sprint(err))
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
