package deer

import "net/http"

func Wrap(w http.ResponseWriter, r *http.Request) (*ResponseWriterAdapter, *RequestAdapter) {
	return WrapResponseWriter(w), WrapRequest(r)
}

type HandlerFunc func(w *ResponseWriterAdapter, r *RequestAdapter)

func (handler HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler(Wrap(w, r))
}

func HandleFunc(pattern string, handlerFunc HandlerFunc) {
	http.Handle(pattern, handlerFunc)
}
