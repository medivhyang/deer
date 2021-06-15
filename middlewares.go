package deer

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

type Middleware = func(HandlerFunc) HandlerFunc

func AllowedMethods(methods ...string) Middleware {
	return func(h HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) {
			find := false
			for _, method := range methods {
				if r.Method() == method {
					find = true
				}
			}
			if !find {
				w.StatusCode(http.StatusMethodNotAllowed)
				return
			}
			h.Next(w, r)
		}
	}
}

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
		panic(newError("basic auth with func", "require func"))
	}
	var finalRealm string
	if len(realm) > 0 {
		finalRealm = realm[0]
	} else {
		finalRealm = ""
	}
	return func(h HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) {
			username, password, ok := r.BasicAuth()
			if ok {
				if f(username, password) {
					h.ServeHTTP(w.Raw(), r.Raw)
					return
				}
			}
			w.Header("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", finalRealm))
			w.StatusCode(http.StatusUnauthorized)
			return
		}
	}
}

type CORSOptions struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
}

func CORS(config ...CORSOptions) Middleware {
	var finalConfig CORSOptions
	if len(config) > 0 {
		finalConfig = config[0]
	} else {
		finalConfig = CORSOptions{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"*"},
			AllowHeaders: []string{"*"},
		}
	}
	return func(h HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) {
			if len(finalConfig.AllowOrigins) > 0 {
				w.Header("Access-Control-Allow-Origin", strings.Join(finalConfig.AllowOrigins, ","))
			}
			if len(finalConfig.AllowMethods) > 0 {
				w.Header("Access-Control-Allow-Methods", strings.Join(finalConfig.AllowMethods, ","))
			}
			if len(finalConfig.AllowHeaders) > 0 {
				w.Header("Access-Control-Allow-Headers", strings.Join(finalConfig.AllowHeaders, ","))
			}
			if len(finalConfig.ExposeHeaders) > 0 {
				w.Header("Access-Control-Expose-Headers", strings.Join(finalConfig.ExposeHeaders, ","))
			}
			if finalConfig.AllowCredentials {
				w.Header("Access-Control-Allow-Credentials", "true")
			}
			h.Next(w, r)
		}
	}
}

func MaxAllowed(n int) Middleware {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }
	return func(f HandlerFunc) HandlerFunc {
		acquire()
		defer release()
		return func(w ResponseWriter, r *Request) {
			f.Next(w, r)
		}
	}
}

func Recovery(callback ...func(w ResponseWriter, r *Request, err interface{})) Middleware {
	var f func(w ResponseWriter, r *Request, err interface{})
	if len(callback) > 0 {
		f = callback[0]
	} else {
		f = func(w ResponseWriter, r *Request, err interface{}) {
			debugf("deer: recovery: %+v\n%s", err, string(debug.Stack()))
			w.Text(http.StatusInternalServerError, fmt.Sprint(err))
		}
	}
	return func(h HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) {
			defer func() {
				if err := recover(); err != nil {
					f(w, r, err)
				}
			}()
			h.ServeHTTP(w.Raw(), r.Raw)
		}
	}
}

func Timing(callback ...func(w ResponseWriter, r *Request, d time.Duration)) Middleware {
	var f func(w ResponseWriter, r *Request, d time.Duration)
	if len(callback) > 0 {
		f = callback[0]
	} else {
		f = func(w ResponseWriter, r *Request, d time.Duration) {
			debugf("timing: \"%s %s\" cost %s\n", r.Method(), r.Path(), d)
		}
	}
	return func(h HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) {
			defer func(start time.Time) {
				d := time.Since(start)
				f(w, r, d)
			}(time.Now())
			h.Next(w, r)
		}
	}
}

func Trace(callback ...func(w ResponseWriter, r *Request)) Middleware {
	var f func(w ResponseWriter, r *Request)
	if len(callback) > 0 {
		f = callback[0]
	} else {
		f = func(w ResponseWriter, r *Request) {
			debugf("%s %s", r.Method(), r.Path())
		}
	}
	return func(h HandlerFunc) HandlerFunc {
		return func(w ResponseWriter, r *Request) {
			f(w, r)
			h.Next(w, r)
		}
	}
}
