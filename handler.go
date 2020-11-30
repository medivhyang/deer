package deer

import "net/http"

type HandlerFunc func(w ResponseWriter, r *Request)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(WrapResponseWriter(w), WrapRequest(r))
}

func (h HandlerFunc) Next(w ResponseWriter, r *Request) {
	h.ServeHTTP(w.Raw(), r.Raw())
}

func WrapHandler(h http.Handler) HandlerFunc {
	return func(w ResponseWriter, r *Request) {
		h.ServeHTTP(w.Raw(), r.Raw())
	}
}

func UnwrapHandler(h HandlerFunc) http.Handler {
	return h
}
