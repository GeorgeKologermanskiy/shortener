package main

import (
	"fmt"
	"log"
	"net/http"

    "github.com/gorilla/mux"
)

func createRouter(storageWriter *StorageWriter) (r *mux.Router) {
	r = mux.NewRouter()
	r.Handle("/create_link", CreateLinkHandler { storageWriter }).Methods(http.MethodPost)
	r.Handle("/delete_link", DeleteLinkHandler { storageWriter }).Methods(http.MethodDelete)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	}).Methods("GET")
	return
}

func main() {
	log.Print("Starting service...")

	log.Print("Creating storage writer...")
	// TODO: move to other servide...
	storageWriter := newStorageWriter("http://localhost:8080")

	log.Print("Creating router with handlers...")
	r := createRouter(storageWriter)

	log.Print("Start serving...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}