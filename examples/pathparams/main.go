package main

import (
	"fmt"
	"net/http"

	"github.com/medivhyang/deer"
	"github.com/medivhyang/deer/middlewares"
)

func main() {
	d := deer.NewRouter().Use(middlewares.Trace())

	d.Get("/orgs/:oid", func(w *deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, fmt.Sprintf("oid = %s", r.PathParam("oid")))
	})
	d.Get("/orgs/:oid/users/:uid", func(w *deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, fmt.Sprintf("oid = %s, uid = %s", r.PathParam("oid"), r.PathParam("uid")))
	})
	d.Get("/static/*filename", func(w *deer.ResponseWriter, r *deer.Request) {
		w.Text(http.StatusOK, fmt.Sprintf("filename = %s", r.PathParam("filename")))
	})

	d.Run(":8080")
}
