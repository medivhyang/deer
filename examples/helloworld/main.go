package main

import (
	"net/http"

	"github.com/medivhyang/deer"
	"github.com/medivhyang/deer/middlewares"
)

func main() {
	d := deer.New().Use(middlewares.Trace())

	d.Get("/", deer.HandlerFunc(func(w *deer.ResponseWriterAdapter, r *deer.RequestAdapter) {
		w.Text(http.StatusOK, "hello world")
	}))

	d.Run(":8080")
}