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
