package deer

import "net/http"

type HandlerFunc func(w *ResponseWriter, r *Request)

func (handler HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler(WrapResponseWriter(w), WrapRequest(r))
}
