package main

import (
	"net/http"

	"github.com/medivhyang/deer"
	"github.com/medivhyang/deer/middlewares"
)

func main() {
	d := deer.NewRouter().Use(middlewares.Trace())

	d.Get("/", func(w *deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, "hello world")
	})

	d.Run(":8080")
}
