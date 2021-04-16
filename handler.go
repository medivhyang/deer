package deer

import "net/http"

type HandlerFunc func(w ResponseWriter, r *Request)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(WrapResponseWriter(w), WrapRequest(r))
}

func (h HandlerFunc) Next(w ResponseWriter, r *Request) {
	h.ServeHTTP(w.Raw(), r.Raw())
}
