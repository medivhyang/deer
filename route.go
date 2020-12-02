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

type Router struct {
	prefix          string
	entryMap        map[key]*entry
	entries         []*entry
	middlewares     []Middleware
	notFoundHandler HandlerFunc
}

type key struct {
	method string
	path   string
}

type entry struct {
	method      string
	pattern     string
	handler     HandlerFunc
	middlewares []Middleware
	regexp      *regexp.Regexp
}

func (e *entry) regexpMatch(path string) bool {
	if e.regexp == nil {
		e.regexp = regexp.MustCompile(toRegexp(e.pattern))
	}
	return e.regexp.MatchString(path)
}

func (e *entry) params(path string) map[string]string {
	if e.regexp == nil {
		e.regexp = regexp.MustCompile(toRegexp(e.pattern))
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

var (
	namedParamRegexp    = regexp.MustCompile(":([^/]+)")
	wildcardParamRegexp = regexp.MustCompile("\\*([^/]+)")
)

func toRegexp(pattern string) string {
	s := namedParamRegexp.ReplaceAllString(pattern, "(?P<$1>[^/]+)")
	s = wildcardParamRegexp.ReplaceAllString(s, "(?P<$1>.*)")
	s = "^" + s + "$"
	return s
}

func NewRouter() *Router {
	return &Router{}
}

func (router *Router) Prefix(p string) *Router {
	router.prefix = normalizePrefix(p)
	return router
}

func (router *Router) HandleNotFound(h HandlerFunc) *Router {
	router.notFoundHandler = h
	return router
}

func (router *Router) Use(middlewares ...Middleware) *Router {
	router.middlewares = append(router.middlewares, middlewares...)
	return router
}

func (router *Router) Handle(method string, path string, handler HandlerFunc, middlewares ...Middleware) *Router {
	router.handle(method, path, handler, middlewares...)
	return router
}

func (router *Router) handle(method string, path string, handler HandlerFunc, middlewares ...Middleware) *Router {
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
			if v.regexpMatch(path) && (v.method == "" || v.method == method) {
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
		r = r.WithContext(context.WithValue(r.Context(), paramsContextKeySingleton, e.params(path)))
	}
	h := Chain(Chain(e.handler, e.middlewares...), router.middlewares...)
	h.ServeHTTP(w, r)
}

func (router *Router) Group(prefix string, middlewares ...Middleware) *group {
	return &group{router: router, prefix: normalizePrefix(prefix), middlewares: middlewares}
}

func (router *Router) Any(pattern string, handler HandlerFunc, middlewares ...Middleware) *Router {
	router.Handle(http.MethodGet, pattern, handler, middlewares...)
	router.Handle(http.MethodPost, pattern, handler, middlewares...)
	router.Handle(http.MethodPut, pattern, handler, middlewares...)
	router.Handle(http.MethodPatch, pattern, handler, middlewares...)
	router.Handle(http.MethodDelete, pattern, handler, middlewares...)
	router.Handle(http.MethodOptions, pattern, handler, middlewares...)
	return router
}

func (router *Router) Get(pattern string, handler HandlerFunc, middlewares ...Middleware) *Router {
	router.Handle(http.MethodGet, pattern, handler, middlewares...)
	router.Handle(http.MethodOptions, pattern, handler, middlewares...)
	return router
}

func (router *Router) Post(pattern string, handler HandlerFunc, middlewares ...Middleware) *Router {
	router.Handle(http.MethodPost, pattern, handler, middlewares...)
	router.Handle(http.MethodOptions, pattern, handler, middlewares...)
	return router
}

func (router *Router) Put(pattern string, handler HandlerFunc, middlewares ...Middleware) *Router {
	router.Handle(http.MethodPut, pattern, handler, middlewares...)
	router.Handle(http.MethodOptions, pattern, handler, middlewares...)
	return router
}

func (router *Router) Patch(pattern string, handler HandlerFunc, middlewares ...Middleware) *Router {
	router.Handle(http.MethodPatch, pattern, handler, middlewares...)
	router.Handle(http.MethodOptions, pattern, handler, middlewares...)
	return router
}

func (router *Router) Delete(pattern string, handler HandlerFunc, middlewares ...Middleware) *Router {
	router.Handle(http.MethodDelete, pattern, handler, middlewares...)
	router.Handle(http.MethodOptions, pattern, handler, middlewares...)
	return router
}

func (router *Router) Options(pattern string, handler HandlerFunc, middlewares ...Middleware) *Router {
	return router.Handle(http.MethodOptions, pattern, handler, middlewares...)
}

type View struct {
	Method  string `json:"method"`
	Pattern string `json:"pattern"`
}

type sortedEntrySlice []*entry

func (s sortedEntrySlice) Len() int {
	return len(s)
}

var sortMethods = map[string]int{
	"":                 0,
	http.MethodGet:     1,
	http.MethodHead:    2,
	http.MethodPost:    3,
	http.MethodPut:     4,
	http.MethodPatch:   5,
	http.MethodDelete:  6,
	http.MethodConnect: 7,
	http.MethodOptions: 8,
	http.MethodTrace:   9,
}

func (s sortedEntrySlice) Less(i, j int) bool {
	if s[i].pattern != s[j].pattern {
		return s[i].pattern < s[j].pattern
	}
	mi := strings.ToUpper(s[i].method)
	mj := strings.ToUpper(s[j].method)
	return sortMethods[mi] < sortMethods[mj]
}

func (s sortedEntrySlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (router *Router) Items() []View {
	copies := append([]*entry{}, router.entries...)
	sort.Sort(sortedEntrySlice(copies))
	var result []View
	for _, entry := range copies {
		method, pattern := entry.method, entry.pattern
		if method == "" {
			method = "ANY"
		}
		pattern = router.prefix + pattern
		result = append(result, View{
			Method:  method,
			Pattern: pattern,
		})
	}
	return result
}

func (router *Router) String() string {
	builder := bytes.Buffer{}
	items := router.Items()
	for _, item := range items {
		builder.WriteString(fmt.Sprintf("%-7s %s\n", item.Method, item.Pattern))
	}
	return strings.TrimSuffix(builder.String(), "\n")
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

func Chain(h HandlerFunc, middlewares ...Middleware) HandlerFunc {
	for _, m := range middlewares {
		if m != nil {
			h = WrapHandler(http.HandlerFunc(m(h).ServeHTTP))
		}
	}
	return h
}

type group struct {
	router      *Router
	prefix      string
	middlewares []Middleware
}

func (g *group) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.router.ServeHTTP(w, r)
}

func (g *group) Handle(method string, path string, handler HandlerFunc, middlewares ...Middleware) *group {
	path = g.prefix + path
	finalMiddlewares := append([]Middleware{}, middlewares...)
	finalMiddlewares = append(finalMiddlewares, g.middlewares...)
	g.router.Handle(method, path, handler, finalMiddlewares...)
	return g
}

func (g *group) Group(prefix string) *group {
	return &group{router: g.router, prefix: g.prefix + normalizePrefix(prefix)}
}

func (g *group) Use(middlewares ...Middleware) *group {
	g.middlewares = append(g.middlewares, middlewares...)
	return g
}

func (g *group) Any(pattern string, handler HandlerFunc, middlewares ...Middleware) *group {
	g.Handle(http.MethodGet, pattern, handler, middlewares...)
	g.Handle(http.MethodPost, pattern, handler, middlewares...)
	g.Handle(http.MethodPut, pattern, handler, middlewares...)
	g.Handle(http.MethodPatch, pattern, handler, middlewares...)
	g.Handle(http.MethodDelete, pattern, handler, middlewares...)
	g.Handle(http.MethodOptions, pattern, handler, middlewares...)
	return g
}

func (g *group) Get(pattern string, handler HandlerFunc, middlewares ...Middleware) *group {
	g.Handle(http.MethodGet, pattern, handler, middlewares...)
	g.Handle(http.MethodOptions, pattern, handler, middlewares...)
	return g
}

func (g *group) Post(pattern string, handler HandlerFunc, middlewares ...Middleware) *group {
	g.Handle(http.MethodPost, pattern, handler, middlewares...)
	g.Handle(http.MethodOptions, pattern, handler, middlewares...)
	return g
}

func (g *group) Put(pattern string, handler HandlerFunc, middlewares ...Middleware) *group {
	g.Handle(http.MethodPut, pattern, handler, middlewares...)
	g.Handle(http.MethodOptions, pattern, handler, middlewares...)
	return g
}

func (g *group) Patch(pattern string, handler HandlerFunc, middlewares ...Middleware) *group {
	g.Handle(http.MethodPatch, pattern, handler, middlewares...)
	g.Handle(http.MethodOptions, pattern, handler, middlewares...)
	return g
}

func (g *group) Delete(pattern string, handler HandlerFunc, middlewares ...Middleware) *group {
	g.Handle(http.MethodDelete, pattern, handler, middlewares...)
	g.Handle(http.MethodOptions, pattern, handler, middlewares...)
	return g
}

func (g *group) Options(pattern string, handler HandlerFunc, middlewares ...Middleware) *group {
	return g.Handle(http.MethodOptions, pattern, handler, middlewares...)
}

type paramsContextKey struct{}

var paramsContextKeySingleton = paramsContextKey{}

func Params(r *http.Request) map[string]string {
	m, ok := r.Context().Value(paramsContextKeySingleton).(map[string]string)
	if ok {
		return m
	}
	return map[string]string{}
}
