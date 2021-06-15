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

// region router

type (
	Router struct {
		prefix          string
		entryMap        map[key]*entry
		entries         []*entry
		middlewares     []Middleware
		notFoundHandler HandlerFunc
	}
	key struct {
		method string
		path   string
	}
	entry struct {
		method      string
		pattern     string
		handler     HandlerFunc
		middlewares []Middleware
		regexp      *regexp.Regexp
	}
)

func (e *entry) regexpMatch(path string) bool {
	if e.regexp == nil {
		e.regexp = regexp.MustCompile(parsePatternToRegexp(e.pattern))
	}
	return e.regexp.MatchString(path)
}

func (e *entry) params(path string) map[string]string {
	if e.regexp == nil {
		e.regexp = regexp.MustCompile(parsePatternToRegexp(e.pattern))
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
	var e *entry
	if strings.HasPrefix(path, router.prefix) {
		path = strings.TrimPrefix(path, router.prefix)
		k := key{
			method: method,
			path:   path,
		}
		if router.entryMap == nil {
			router.entryMap = map[key]*entry{}
		}
		e = router.entryMap[k]
		if e == nil {
			regexpMatch := false
			for _, v := range router.entries {
				if v.regexpMatch(path) && v.method == method {
					regexpMatch = true
					e = v
					break
				}
			}
			if regexpMatch && e != nil {
				r = r.WithContext(context.WithValue(r.Context(), paramsContextKeySingleton, e.params(path)))
			}
		}
	}
	if e == nil {
		if method == http.MethodOptions {
			e = &entry{
				method:  http.MethodOptions,
				handler: defaultOptionsHandleFunc,
			}
		} else {
			e = &entry{
				method:  method,
				handler: defaultNotFoundHandleFunc,
			}
		}
	}
	finalMiddlewares := append([]Middleware{}, router.middlewares...)
	finalMiddlewares = append(finalMiddlewares, e.middlewares...)
	h := chain(e.handler, finalMiddlewares...)
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
	return router
}

func (router *Router) Post(pattern string, handler HandlerFunc, middlewares ...Middleware) *Router {
	router.Handle(http.MethodPost, pattern, handler, middlewares...)
	return router
}

func (router *Router) Put(pattern string, handler HandlerFunc, middlewares ...Middleware) *Router {
	router.Handle(http.MethodPut, pattern, handler, middlewares...)
	return router
}

func (router *Router) Patch(pattern string, handler HandlerFunc, middlewares ...Middleware) *Router {
	router.Handle(http.MethodPatch, pattern, handler, middlewares...)
	return router
}

func (router *Router) Delete(pattern string, handler HandlerFunc, middlewares ...Middleware) *Router {
	router.Handle(http.MethodDelete, pattern, handler, middlewares...)
	return router
}

func (router *Router) Options(pattern string, handler HandlerFunc, middlewares ...Middleware) *Router {
	return router.Handle(http.MethodOptions, pattern, handler, middlewares...)
}

type EntryView struct {
	Method  string `json:"method"`
	Pattern string `json:"pattern"`
}

var methodOrders = map[string]int{
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

func (router *Router) Items() []EntryView {
	items := append([]*entry{}, router.entries...)
	sort.Slice(items, func(i, j int) bool {
		if items[i].pattern != items[j].pattern {
			return items[i].pattern < items[j].pattern
		}
		mi := strings.ToUpper(items[i].method)
		mj := strings.ToUpper(items[j].method)
		return methodOrders[mi] < methodOrders[mj]
	})
	var result []EntryView
	for _, item := range items {
		method, pattern := item.method, item.pattern
		if method == "" {
			method = "ANY"
		}
		pattern = router.prefix + pattern
		result = append(result, EntryView{
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

// endregion

// region router group

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
	finalMiddlewares := append([]Middleware{}, g.middlewares...)
	finalMiddlewares = append(finalMiddlewares, middlewares...)
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
	return g
}

func (g *group) Post(pattern string, handler HandlerFunc, middlewares ...Middleware) *group {
	g.Handle(http.MethodPost, pattern, handler, middlewares...)
	return g
}

func (g *group) Put(pattern string, handler HandlerFunc, middlewares ...Middleware) *group {
	g.Handle(http.MethodPut, pattern, handler, middlewares...)
	return g
}

func (g *group) Patch(pattern string, handler HandlerFunc, middlewares ...Middleware) *group {
	g.Handle(http.MethodPatch, pattern, handler, middlewares...)
	return g
}

func (g *group) Delete(pattern string, handler HandlerFunc, middlewares ...Middleware) *group {
	g.Handle(http.MethodDelete, pattern, handler, middlewares...)
	return g
}

func (g *group) Options(pattern string, handler HandlerFunc, middlewares ...Middleware) *group {
	return g.Handle(http.MethodOptions, pattern, handler, middlewares...)
}

// endregion

// region params

type paramsContextKey struct{}

var paramsContextKeySingleton = paramsContextKey{}

func Params(r *http.Request) map[string]string {
	m, ok := r.Context().Value(paramsContextKeySingleton).(map[string]string)
	if ok {
		return m
	}
	return map[string]string{}
}

// endregion

// region utils

var (
	namedParamRegexp    = regexp.MustCompile(":([^/]+)")
	wildcardParamRegexp = regexp.MustCompile("\\*([^/]+)")
)

func parsePatternToRegexp(pattern string) string {
	s := namedParamRegexp.ReplaceAllString(pattern, "(?P<$1>[^/]+)")
	s = wildcardParamRegexp.ReplaceAllString(s, "(?P<$1>.*)")
	s = "^" + s + "$"
	return s
}

func defaultOptionsHandleFunc(w ResponseWriter, r *Request) {
	w.Header("Content-Length", "0")
	w.StatusCode(http.StatusNoContent)
}

func defaultNotFoundHandleFunc(w ResponseWriter, r *Request) {
	w.Text(http.StatusNotFound, http.StatusText(http.StatusNotFound))
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
	findIndex := sort.Search(n, func(i int) bool {
		return es[i].method == e.method && es[i].pattern == e.pattern
	})
	if findIndex < n {
		es[findIndex] = e
		return es
	}
	smallestIndex := sort.Search(n, func(i int) bool {
		l1 := len(strings.Split(es[i].pattern, "/"))
		l2 := len(strings.Split(e.pattern, "/"))
		if l1 != l2 {
			return l1 < l2
		}
		return len(es[i].pattern) < len(e.pattern)
	})
	if smallestIndex == n {
		return append(es, e)
	}
	es = append(es, nil)
	copy(es[smallestIndex+1:], es[smallestIndex:])
	es[smallestIndex] = e
	return es
}

func chain(h HandlerFunc, middlewares ...Middleware) HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// endregion
