package main

import (
	"log"
	"net/http"

	"github.com/medivhyang/deer"
)

func main() {
	deer.Default().Use(deer.Recovery(func(w deer.ResponseWriter, r *deer.Request, err interface{}) {

	}))
	r := deer.Default()

	r.Get("/", func(w deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, "hello world")
	})

	log.Fatalln(r.Run(":8080"))
}
