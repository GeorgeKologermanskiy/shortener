package main

import (
	"log"
	"net/http"

    "github.com/gorilla/mux"
)

func createRouter(storageWriter *StorageWriter) (r *mux.Router) {
	r = mux.NewRouter()
	s := r.PathPrefix("/api/manage_links").Subrouter()
	s.Handle("/create_link", CreateLinkHandler { storageWriter }).Methods(http.MethodPost)
	s.Handle("/delete_link", DeleteLinkHandler { storageWriter }).Methods(http.MethodDelete)
	s.Handle("/links_info", LinksInfoHandler { storageWriter }).Methods(http.MethodGet)
	return
}

func main() {
	log.Print("Starting service...")

	log.Print("Creating storage writer...")
	// TODO: move to other service...
	storageWriter := newStorageWriter("http://localhost:8087")

	log.Print("Creating router with handlers...")
	r := createRouter(storageWriter)

	log.Print("Start serving...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}