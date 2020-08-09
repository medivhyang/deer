package main

import (
	"fmt"
	"net/http"

	"github.com/medivhyang/deer"
	"github.com/medivhyang/deer/middlewares"
)

func main() {
	d := deer.New().Use(middlewares.Trace())

	d.Get("/orgs/:oid", deer.HandlerFunc(func(w *deer.ResponseWriterAdapter, r *deer.RequestAdapter) {
		w.Text(http.StatusOK, fmt.Sprintf("oid = %s", r.PathParam("oid")))
	}))
	d.Get("/orgs/:oid/users/:uid", deer.HandlerFunc(func(w *deer.ResponseWriterAdapter, r *deer.RequestAdapter) {
		w.Text(http.StatusOK, fmt.Sprintf("oid = %s, uid = %s", r.PathParam("oid"), r.PathParam("uid")))
	}))
	d.Get("/static/*filename", deer.HandlerFunc(func(w *deer.ResponseWriterAdapter, r *deer.RequestAdapter) {
		w.Text(http.StatusOK, fmt.Sprintf("filename = %s", r.PathParam("filename")))
	}))

	d.Run(":8080")
}
