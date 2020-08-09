# Deer

Deer is a go http libary, not framework.

## Why Deer

1. Lightweight, no have too much concept.
2. Less intrusion, just libary, not framework.
3. Native, compatible with standard libary.
4. Rich Feature, support group route, path params, and middleware etc.
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
	d := deer.New().Use(middlewares.Trace())

	d.Get("/", deer.HandlerFunc(func(w *deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, "hello world")
	}))

	d.Run(":8080")
}
```

Demo path params usage:

```go
func main() {
	d := deer.New().Use(middlewares.Trace())

	d.Get("/orgs/:oid", deer.HandlerFunc(func(w *deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, fmt.Sprintf("oid = %s", r.PathParam("oid")))
	}))
	d.Get("/orgs/:oid/users/:uid", deer.HandlerFunc(func(w *deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, fmt.Sprintf("oid = %s, uid = %s", r.PathParam("oid"), r.PathParam("uid")))
	}))
	d.Get("/static/*filename", deer.HandlerFunc(func(w *deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, fmt.Sprintf("filename = %s", r.PathParam("filename")))
	}))

	d.Run(":8080")
}
```

Demo http client usage:

```go
func main() {
	resp, err := deer.Get("https://example.com").Do()
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Text())
}
```

More examples references `/examples` directory.