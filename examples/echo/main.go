package main

import (
	"fmt"
	"net/http"

	"github.com/medivhyang/deer"
)

func main() {
	d := deer.New()

	d.Any("/", deer.HandlerFunc(func(w *deer.ResponseWriterAdapter, r *deer.RequestAdapter) {
		w.Text(http.StatusOK, "welcome to echo!")
	}))
	d.Any("/echo", deer.HandlerFunc(func(w *deer.ResponseWriterAdapter, r *deer.RequestAdapter) {
		w.Text(http.StatusOK, "echo: "+r.Query("q"))
	}))

	g := d.Group("v2")
	{
		g.Any("/echo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("echo v2: " + r.FormValue("q")))
		}))
	}

	fmt.Println(d)

	d.Run(":8080")
}
