package deer

import "net/http"

type HandlerFunc func(w *ResponseWriterAdapter, r *RequestAdapter)

func (handler HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler(WrapResponseWriter(w), WrapRequest(r))
}
