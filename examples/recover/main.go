package main

import (
	"github.com/medivhyang/deer"
	"github.com/medivhyang/deer/middlewares"
)

func main() {
	r := deer.NewRouter().Use(middlewares.Trace(), middlewares.Recovery())

	r.Get("/", func(w deer.ResponseWriter, r *deer.Request) {
		panic("1")
	})

	r.Run(":8080")
}
