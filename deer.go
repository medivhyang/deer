package deer

import (
	"fmt"
	"io"
	"os"
)

var (
	debugFlag             = false
	debugWriter io.Writer = os.Stdout
)

func Debug(b bool) {
	debugFlag = b
}

func Output(writer io.Writer) {
	debugWriter = writer
}

func Default() *Router {
	return NewRouter().Use(Recovery(), Trace())
}

func debugf(format string, args ...interface{}) {
	if !debugFlag {
		return
	}
	if debugWriter == nil {
		return
	}
	if _, err := fmt.Fprintln(debugWriter, fmt.Sprintf(format, args...)); err != nil {
		panic(err)
	}
}

func newError(module string, format string, args ...interface{}) error {
	return fmt.Errorf("%s: %s: %s", "deer", module, fmt.Sprintf(format, args...))
}
