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
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Enabling Cors to write Vue App consuming the REST API
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// FileServer may be used to serve static files from a directory
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := chi.NewRouter()
	r.Use(commonMiddleware)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "/")
	FileServer(r, "/index.html", http.Dir(filesDir))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		w.Write([]byte("."))
	})
	// Starts the worker and releases the Kraken
	r.Get("/worker", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range database {
			go polling(v.INTERVAL*1000000000, v.URL, k)
		}
		w.WriteHeader(200)
		w.Write([]byte("..."))
	})

	r.Mount("/api/fetcher", Work{}.Routes())

	http.ListenAndServe(":8080", r)
}
