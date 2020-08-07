package deer

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type Middleware = func(http.Handler) http.Handler

type Router struct {
	prefix          string
	entryMap        map[key]*entry
	entries         []*entry
	middlewares     []Middleware
	notFoundHandler http.Handler
}

type key struct {
	method string
	path   string
}

type entry struct {
	method      string
	pattern     string
	handler     http.Handler
	middlewares []Middleware
	regexp      *regexp.Regexp
}

func (e *entry) regexpMatch(path string) bool {
	if e.regexp == nil {
		e.regexp = regexp.MustCompile(trPatternToRegexp(e.pattern))
	}
	return e.regexp.MatchString(path)
}

func (e *entry) pathParams(path string) map[string]string {
	if e.regexp == nil {
		e.regexp = regexp.MustCompile(trPatternToRegexp(e.pattern))
	}
	matches := e.regexp.FindStringSubmatch(path)
	names := e.regexp.SubexpNames()
	result := map[string]string{}
	for i, name := range names {
		if i > 0 {
			result[name] = matches[i]
		}
	}
	return result
}

var namedParamPattern = regexp.MustCompile(":([^/]+)")
var wildcardParamPattern = regexp.MustCompile("\\*([^/]+)")

func trPatternToRegexp(pattern string) string {
	s := namedParamPattern.ReplaceAllString(pattern, "(?P<$1>[^/]+)")
	s = wildcardParamPattern.ReplaceAllString(s, "(?P<$1>.*)")
	s = "^" + s + "$"
	return s
}

func New() *Router {
	return &Router{}
}

func (router *Router) Prefix(p string) *Router {
	router.prefix = normalizePrefix(p)
	return router
}

func (router *Router) HandleNotFound(h http.Handler) *Router {
	router.notFoundHandler = h
	return router
}

func (router *Router) Use(middlewares ...Middleware) *Router {
	router.middlewares = append(router.middlewares, middlewares...)
	return router
}

func (router *Router) Handle(method string, path string, handler http.Handler, middlewares ...Middleware) *Router {
	path = normalizePath(path)
	e := entry{
		method:      method,
		pattern:     path,
		handler:     handler,
		middlewares: middlewares,
	}
	k := key{
		method: method,
		path:   path,
	}
	if router.entryMap == nil {
		router.entryMap = map[key]*entry{}
	}
	if _, ok := router.entryMap[k]; ok {
		panic(fmt.Sprintf("deer: route \"%s %s\" has registered", method, method))
	}
	router.entryMap[k] = &e
	router.entries = appendSorted(router.entries, &e)
	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method, path := r.Method, r.URL.Path
	if !strings.HasPrefix(path, router.prefix) {
		router.notFound(w, r)
		return
	}
	path = strings.TrimPrefix(path, router.prefix)
	k := key{
		method: method,
		path:   path,
	}
	if router.entryMap == nil {
		router.entryMap = map[key]*entry{}
	}
	e := router.entryMap[k]
	if e == nil {
		k.method = ""
		e = router.entryMap[k]
	}
	regexpMatch := false
	if e == nil {
		for _, v := range router.entries {
			if v.regexpMatch(path) {
				regexpMatch = true
				e = v
				break
			}
		}
	}
	if e == nil || e.handler == nil {
		router.notFound(w, r)
		return
	}
	if regexpMatch {
		r = r.WithContext(context.WithValue(r.Context(), pathParamsContextKey, e.pathParams(path)))
	}
	h := chain(chain(e.handler, e.middlewares...), router.middlewares...)
	h.ServeHTTP(w, r)
}

func (router *Router) Group(prefix string) *group {
	return &group{router: router, prefix: normalizePrefix(prefix)}
}

func (router *Router) Any(pattern string, handler http.Handler, middlewares ...Middleware) *Router {
	return router.Handle("", pattern, handler, middlewares...)
}

func (router *Router) Get(pattern string, handler http.Handler, middlewares ...Middleware) *Router {
	return router.Handle(http.MethodGet, pattern, handler, middlewares...)
}

func (router *Router) Post(pattern string, handler http.Handler, middlewares ...Middleware) *Router {
	return router.Handle(http.MethodPost, pattern, handler, middlewares...)
}

func (router *Router) Put(pattern string, handler http.Handler, middlewares ...Middleware) *Router {
	return router.Handle(http.MethodPut, pattern, handler, middlewares...)
}

func (router *Router) Patch(pattern string, handler http.Handler, middlewares ...Middleware) *Router {
	return router.Handle(http.MethodPatch, pattern, handler, middlewares...)
}

func (router *Router) Delete(pattern string, handler http.Handler, middlewares ...Middleware) *Router {
	return router.Handle(http.MethodDelete, pattern, handler, middlewares...)
}

func (router *Router) Options(pattern string, handler http.Handler, middlewares ...Middleware) *Router {
	return router.Handle(http.MethodOptions, pattern, handler, middlewares...)
}

func (router *Router) String() string {
	builder := bytes.Buffer{}
	builder.WriteString("---\n")

	for i := len(router.entries) - 1; i >= 0; i-- {
		method, path := router.entries[i].method, router.entries[i].pattern
		if method == "" {
			method = "ANY"
		}
		path = router.prefix + path
		builder.WriteString(fmt.Sprintf("%s %s\n", method, path))
	}
	builder.WriteString("---\n")
	return builder.String()
}

func (router *Router) Run(addr string) error {
	return http.ListenAndServe(addr, router)
}

func (router *Router) notFound(w http.ResponseWriter, r *http.Request) {
	if router.notFoundHandler != nil {
		router.notFoundHandler.ServeHTTP(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(http.StatusText(http.StatusNotFound)))
}

func normalizePath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		return "/" + p
	}
	return p
}

func normalizePrefix(p string) string {
	p = strings.TrimRight(p, "/")
	if p == "" {
		return ""
	}
	if p[0] != '/' {
		return "/" + p
	}
	return p
}

func appendSorted(es []*entry, e *entry) []*entry {
	n := len(es)
	i := sort.Search(n, func(i int) bool {
		l1 := len(strings.Split(es[i].pattern, "/"))
		l2 := len(strings.Split(e.pattern, "/"))
		if l1 != l2 {
			return l1 < l2
		}
		return len(es[i].pattern) < len(e.pattern)
	})
	if i == n {
		return append(es, e)
	}
	es = append(es, nil)
	copy(es[i+1:], es[i:])
	es[i] = e
	return es
}

func chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		if m != nil {
			h = m(h)
		}
	}
	return h
}

type group struct {
	router *Router
	prefix string
}

func (g *group) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.router.ServeHTTP(w, r)
}

func (g *group) Handle(method string, path string, handler http.Handler, middlewares ...Middleware) *group {
	path = g.prefix + path
	g.router.Handle(method, path, handler, middlewares...)
	return g
}

func (g *group) Group(prefix string) *group {
	return &group{router: g.router, prefix: g.prefix + normalizePrefix(prefix)}
}

func (g *group) Any(pattern string, handler http.Handler, middlewares ...Middleware) *group {
	return g.Handle("", pattern, handler, middlewares...)
}

func (g *group) Get(pattern string, handler http.Handler, middlewares ...Middleware) *group {
	return g.Handle(http.MethodGet, pattern, handler, middlewares...)
}

func (g *group) Post(pattern string, handler http.Handler, middlewares ...Middleware) *group {
	return g.Handle(http.MethodPost, pattern, handler, middlewares...)
}

func (g *group) Put(pattern string, handler http.Handler, middlewares ...Middleware) *group {
	return g.Handle(http.MethodPut, pattern, handler, middlewares...)
}

func (g *group) Patch(pattern string, handler http.Handler, middlewares ...Middleware) *group {
	return g.Handle(http.MethodPatch, pattern, handler, middlewares...)
}

func (g *group) Delete(pattern string, handler http.Handler, middlewares ...Middleware) *group {
	return g.Handle(http.MethodDelete, pattern, handler, middlewares...)
}

func (g *group) Options(pattern string, handler http.Handler, middlewares ...Middleware) *group {
	return g.Handle(http.MethodOptions, pattern, handler, middlewares...)
}

var pathParamsContextKey = &struct{}{}

func PathParams(r *http.Request) map[string]string {
	m, ok := r.Context().Value(pathParamsContextKey).(map[string]string)
	if ok {
		return m
	}
	return map[string]string{}
}
