package middlewares

import (
	"fmt"
	"github.com/medivhyang/deer"
	"net/http"
)

func BasicAuth(pairs map[string]string, realm ...string) deer.Middleware {
	return BasicAuthWithFunc(func(username, password string) bool {
		for k, v := range pairs {
			if username == k && password == v {
				return true
			}
		}
		return false
	}, realm...)
}

func BasicAuthWithFunc(f func(username, password string) bool, realm ...string) deer.Middleware {
	if f == nil {
		panic("basic auth with func: require func")
	}
	var finalRealm string
	if len(realm) > 0 {
		finalRealm = realm[0]
	} else {
		finalRealm = ""
	}
	return func(h http.Handler) http.Handler {
		return deer.HandlerFunc(func(w *deer.ResponseWriter, r *deer.Request) {
			username, password, ok := r.Raw.BasicAuth()
			if ok {
				if f(username, password) {
					h.ServeHTTP(w.Raw, r.Raw)
					return
				}
			}
			w.SetHeader("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", finalRealm))
			w.SetStatusCode(http.StatusUnauthorized)
			return
		})
	}
}
