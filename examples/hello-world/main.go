package main

import (
	"log"
	"net/http"

	"github.com/medivhyang/deer"
)

func main() {
	r := deer.Default()

	r.Get("/", func(w deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, "hello world")
	})

	log.Fatalln(r.Run(":8080"))
}
