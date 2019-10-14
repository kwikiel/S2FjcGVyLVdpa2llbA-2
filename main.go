//
// Todos Resource
// ==============
// This example demonstrates a project structure that defines a subrouter and its
// handlers on a struct, and mounting them as subrouters to a parent router.
// See also _examples/rest for an in-depth example of a REST service, and apply
// those same patterns to this structure.
//
package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})
	// Starts the worker and releases the Kraken
	r.Get("/worker", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range database {
			fmt.Printf("key[%v] value[%v]\n", k, v)
			go polling(v.INTERVAL*1000000000, v.URL, k)
		}
		w.WriteHeader(200)
		w.Write([]byte("..."))
	})

	r.Mount("/api/fetcher", Work{}.Routes())

	http.ListenAndServe(":8080", r)
}
