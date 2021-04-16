package deer

import (
	"fmt"
	"io"
)

var (
	LogWriter io.Writer
	Debug     = false
)

func logf(format string, args ...interface{}) {
	if LogWriter == nil {
		return
	}
	if _, err := fmt.Fprintf(LogWriter, format, args...); err != nil {
		panic(err)
	}
}

func debugf(format string, args ...interface{}) {
	if Debug {
		logf(format, args...)
	}
}
