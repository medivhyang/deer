package middlewares

import (
	"github.com/medivhyang/deer"
	"net/http"
)

func AllowedMethods(methods ...string) func(deer.HandlerFunc) deer.HandlerFunc {
	return func(h deer.HandlerFunc) deer.HandlerFunc {
		return func(w deer.ResponseWriter, r *deer.Request) {
			find := false
			for _, method := range methods {
				if r.Method() == method {
					find = true
				}
			}
			if !find {
				w.Error(http.StatusMethodNotAllowed)
				return
			}
			h.Next(w, r)
		}
	}
}