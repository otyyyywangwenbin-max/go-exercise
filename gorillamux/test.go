package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

func main() {
	var dir string
	flag.StringVar(&dir, "dir", "/z_tmp/", "the directory to serve files from. Defaults to temp dir")
	flag.Parse()

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/index", indexHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
