package main

import (
	"github.com/medivhyang/deer"
	"github.com/medivhyang/deer/middlewares"
)

func main() {
	d := deer.NewRouter().Use(middlewares.Trace(), middlewares.Recovery())

	d.Get("/", func(w deer.ResponseWriter, r *deer.Request) {
		panic("1")
	})

	d.Run(":8080")
}
