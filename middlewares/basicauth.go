package middlewares

import (
	"fmt"
	"net/http"
)

func BasicAuth(pairs map[string]string, realm ...string) Middleware {
	return BasicAuthWithFunc(func(username, password string) bool {
		for k, v := range pairs {
			if username == k && password == v {
				return true
			}
		}
		return false
	}, realm...)
}

func BasicAuthWithFunc(f func(username, password string) bool, realm ...string) Middleware {
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
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if ok {
				if f(username, password) {
					h.ServeHTTP(w, r)
					return
				}
			}
			w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", finalRealm))
			w.WriteHeader(http.StatusUnauthorized)
			return
		})
	}
}
