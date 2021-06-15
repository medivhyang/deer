package main

import (
	"log"

	"github.com/medivhyang/deer"
)

func main() {
	r := deer.NewRouter().Use(deer.Trace(), deer.Recovery())

	r.Get("/", func(w deer.ResponseWriter, r *deer.Request) {
		panic("1")
	})

	log.Fatalln(r.Run(":8080"))
}
