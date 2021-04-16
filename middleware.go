package deer

type Middleware = func(HandlerFunc) HandlerFunc
