package deer

import "net/http"

type HandlerFunc func(w ResponseWriter, r *Request)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(WrapResponseWriter(w), WrapRequest(r))
}

func (h HandlerFunc) Next(w ResponseWriter, r *Request) {
	h.ServeHTTP(w.Raw(), r.raw)
}

func WrapHandlerFunc(h http.HandlerFunc) HandlerFunc {
	return func(w ResponseWriter, r *Request) {
		h.ServeHTTP(w.Raw(), r.raw)
	}
}

func UnwrapHandlerFunc(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

func WrapHandler(h http.Handler) HandlerFunc {
	return func(w ResponseWriter, r *Request) {
		h.ServeHTTP(w.Raw(), r.raw)
	}
}

func UnwrapHandler(h HandlerFunc) http.Handler {
	return h
}
