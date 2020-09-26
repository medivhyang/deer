package deer

import "net/http"

type HandlerFunc func(w *ResponseWriter, r *Request)

func (handler HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler(WrapResponseWriter(w), WrapRequest(r))
}

func WrapHandlerFunc(h http.HandlerFunc) HandlerFunc {
	return func(w *ResponseWriter, r *Request) {
		h.ServeHTTP(w.Raw, r.Raw)
	}
}

func UnwrapHandlerFunc(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

func WrapHandler(h http.Handler) HandlerFunc {
	return func(w *ResponseWriter, r *Request) {
		h.ServeHTTP(w.Raw, r.Raw)
	}
}

func UnwrapHandler(h HandlerFunc) http.Handler {
	return h
}
