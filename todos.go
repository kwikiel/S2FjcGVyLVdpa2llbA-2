package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

// Work is unit of Work for fetcher
type Work struct {
	ID       int64
	URL      string
	INTERVAL int64
}

// Download history. Note that time duration is in nanoseconds
type Download struct {
	Response  string
	Duration  time.Duration
	CreatedAt time.Time
}

var database = make(map[int64]Work)
var downloads = make(map[int64][]Download)

// Routes creates a REST router for the todos resource
func (rs Work) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.List)    // GET /todos - read a list of todos
	r.Post("/", rs.Create) // POST /todos - create a new todo and persist it
	r.Put("/", rs.Delete)
	r.Get("/x", rs.Crawler)

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(rs.TodoCtx) // lets have a todos map, and lets actually load/manipulate
		r.Get("/", rs.Get)       // GET /todos/{id} - read a single todo by :id
		r.Put("/", rs.Update)    // PUT /todos/{id} - update a single todo by :id
		r.Delete("/", rs.Delete) // DELETE /todos/{id} - delete a single todo by :id
		r.Get("/history", rs.Sync)
	})

	return r
}

// List all Work to fetch
func (rs Work) List(w http.ResponseWriter, r *http.Request) {
	// Display marshalled json from
	djson, _ := json.Marshal(database)
	w.WriteHeader(200)
	w.Write(djson)
}

// Create Worker fetch resource
func (rs Work) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t Work
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	database[t.ID] = t
	id, _ := json.Marshal(t)
	w.WriteHeader(200)
	w.Write(id)
}

// Get specific work resource
func (rs Work) Get(w http.ResponseWriter, r *http.Request) {
	// Again this is just the idiom for getting the id, little cumbersome
	id := chi.URLParam(r, "id")
	var idx int64
	idx, _ = strconv.ParseInt(id, 10, 64)

	var emp = Work{}

	onetodo := database[idx]
	if onetodo == emp {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 HTTP status code returned!"))

	} else {
		rer, _ := json.Marshal(onetodo)
		w.Write(rer)
	}

}

// Update that overwrites the current
func (rs Work) Update(w http.ResponseWriter, r *http.Request) {
	// Exactly same as Update - maybe response should differ if it's overwrite
	// It's called Upsert
	decoder := json.NewDecoder(r.Body)
	var t Work
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	database[t.ID] = t
	id, _ := json.Marshal(t)
	w.WriteHeader(200)
	w.Write(id)
}

// Delete the fetcher Url, doesn't stop the worker yet
func (rs Work) Delete(w http.ResponseWriter, r *http.Request) {
	// This is just the idiom for getting the id
	// TODO: Research more direct ways to get this
	id := chi.URLParam(r, "id")
	var idx int64
	idx, _ = strconv.ParseInt(id, 10, 64)
	// Replace Todobase with database
	delete(database, idx)
	w.WriteHeader(200)
	w.Write([]byte("Deleted"))
}

// Sync displays download history for specific URL id
func (rs Work) Sync(w http.ResponseWriter, r *http.Request) {
	// It's currently called Sync - it should display Download history
	// Index Context is already passed in
	id := chi.URLParam(r, "id")
	var idx int64
	idx, _ = strconv.ParseInt(id, 10, 64)
	djson, _ := json.Marshal(downloads[idx])
	w.WriteHeader(200)
	w.Write(djson)
}

func fetchurl(url string, id int64) string {
	start := time.Now()
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	tnow := time.Now()
	d1 := Download{Response: string(responseData), Duration: time.Since(start), CreatedAt: tnow}
	downloads[id] = append(downloads[id], d1)
	return "ok"
}

func polling(seconds int64, url string, id int64) {
	for {
		<-time.After(time.Duration(seconds))
		go fetchurl(url, id)
	}
}

// Crawler using goroutines
func (rs Work) Crawler(w http.ResponseWriter, r *http.Request) {

	for k, v := range database {
		fmt.Printf("key[%v] value[%v]\n", k, v)
		go polling(v.INTERVAL*1000000000, v.URL, k)
	}
	w.WriteHeader(200)

	w.Write([]byte("."))
	// Start all working requests
}
