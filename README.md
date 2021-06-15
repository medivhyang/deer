# Deer

A go language http library using decorators.

## Why Deer

Lightweight, lower dependence and efficient.

## Quick Start

```go
package main

import (
	"log"
	"net/http"

	"github.com/medivhyang/deer"
)

func main() {
	deer.Debug(true)

	r := deer.Default()

	r.Get("/", func(w deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, "hello world")
	})

	log.Fatalln(r.Run(":8080"))
}
```

> More examples references `/examples` directory.