package deer

import (
	"fmt"
	"io"
)

var (
	LogWriter io.Writer
	Debug     = false
)

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
