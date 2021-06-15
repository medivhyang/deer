package deer

import (
	"fmt"
	"io"
)

var (
	LogWriter io.Writer
	Debug     = false
)

func Default() *Router {
	return NewRouter().Use(Recovery(), Trace())
}

func debugf(format string, args ...interface{}) {
	if !Debug {
		return
	}
	if LogWriter == nil {
		return
	}
	if _, err := fmt.Fprintf(LogWriter, format, args...); err != nil {
		panic(err)
	}
}

func newError(module string, format string, args ...interface{}) error {
	return fmt.Errorf("%s: %s: %s", "deer", module, fmt.Sprintf(format, args...))
}
