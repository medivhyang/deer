# Deer

A go language http library using decorators.

## Why Deer

Lightweight、lower dependence and efficient.

## Quick Start

Hello World

```go
package main

import (
	"log"
	"net/http"

	"github.com/medivhyang/deer"
	"github.com/medivhyang/deer/middlewares"
)

func main() {
	r := deer.NewRouter().Use(middlewares.Trace())

	r.Get("/", func(w deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, "hello world")
	})

	log.Fatalln(r.Run(":8080"))
}
```

Http Client

```go
package main

import (
	"fmt"
	"github.com/medivhyang/deer"
)

func main() {
	text, err := deer.GetText("https://baidu.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
```

More examples references `/examples` directory.