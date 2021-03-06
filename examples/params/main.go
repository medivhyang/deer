package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/medivhyang/deer"
)

func main() {
	r := deer.Default()

	r.Get("/orgs/:oid", func(w deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, fmt.Sprintf("oid = %s", r.Param("oid")))
	})
	r.Get("/orgs/:oid/users/:uid", func(w deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, fmt.Sprintf("oid = %s, uid = %s", r.Param("oid"), r.Param("uid")))
	})
	r.Get("/static/*filename", func(w deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, fmt.Sprintf("filename = %s", r.Param("filename")))
	})

	log.Fatalln(r.Run(":8080"))
}
