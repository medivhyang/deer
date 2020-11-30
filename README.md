# Deer

Deer is a go http libary, not framework.

## Why Deer

1. Lightweight, not much new concept.
2. Less intrusion, just libary, not framework.
3. Native, compatible with standard libary.
4. Rich features, support group route, path params, and middleware etc.
5. Efficient, follow engineering practice.

## Quick Start

Hello World

```go
package main

import (
	"net/http"

	"github.com/medivhyang/deer"
	"github.com/medivhyang/deer/middlewares"
)

func main() {
	r := deer.NewRouter().Use(middlewares.Trace())

	r.Get("/", func(w deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, "hello world")
	})

	r.Run(":8080")
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