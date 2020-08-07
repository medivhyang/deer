package main

import (
	"fmt"
	"github.com/medivhyang/deer"
	"net/http"
)

func main() {
	d := deer.New()

	d.Get("/orgs/:oid", deer.HandlerFunc(func(w *deer.ResponseWriterAdapter, r *deer.RequestAdapter) {
		w.Text(http.StatusOK, "oid = "+r.PathParam("oid"))
	}))
	d.Get("/orgs/:oid/users/:uid", deer.HandlerFunc(func(w *deer.ResponseWriterAdapter, r *deer.RequestAdapter) {
		w.Text(http.StatusOK, "oid = "+r.PathParam("oid")+", uid = "+r.PathParam("uid"))
	}))
	d.Get("/static/*filename", deer.HandlerFunc(func(w *deer.ResponseWriterAdapter, r *deer.RequestAdapter) {
		w.Text(http.StatusOK, "filename = "+r.PathParam("filename"))
	}))

	fmt.Println(d)

	d.Run(":8081")
}
