package deer

import "net/http"

type Middleware = func(HandlerFunc) HandlerFunc

type StandardMiddleware = func(http.Handler) http.Handler

func WrapMiddleware(middleware StandardMiddleware) Middleware {
	return func(h HandlerFunc) HandlerFunc {
		return WrapHandler(http.HandlerFunc(middleware(h).ServeHTTP))
	}
}

func UnwrapMiddleware(middleware Middleware) StandardMiddleware {
	return func(h http.Handler) http.Handler {
		return middleware(WrapHandler(h))
	}
}
