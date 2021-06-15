package main

import (
	"log"

	"github.com/medivhyang/deer"
)

func main() {
	deer.Debug(true)
	r := deer.NewRouter().Use(deer.Recovery(), deer.Trace())

	r.Get("/", func(w deer.ResponseWriter, r *deer.Request) {
		panic("1")
	})

	log.Fatalln(r.Run(":8080"))
}
