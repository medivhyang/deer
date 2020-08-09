package main

import (
	"github.com/medivhyang/deer"
	"github.com/medivhyang/deer/middlewares"
)

func main() {
	d := deer.New().Use(middlewares.Trace(), middlewares.Recovery())

	d.Get("/", deer.HandlerFunc(func(w *deer.ResponseWriterAdapter, r *deer.RequestAdapter) {
		panic("1")
	}))

	d.Run(":8080")
}
