package middlewares

import (
	"fmt"
	"github.com/medivhyang/deer"
)

func logf(format string, args ...interface{}) {
	if deer.LogWriter == nil {
		return
	}
	if _, err := fmt.Fprintf(deer.LogWriter, format, args...); err != nil {
		panic(err)
	}
}
