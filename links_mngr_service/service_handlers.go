package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type CreateLinkHandler struct {
	storageWriter *StorageWriter
}

func (handler CreateLinkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var link string
	if linkList := r.URL.Query()["link"]; len(linkList) == 1 {
		link = string(linkList[0])
	} else {
		log.Print("Incoming bad request on creation link")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: check link before storing...
	linkInfo := ShortLinkInfo {
		Full_link: link,
		Created_time: time.Now(),
		Expired_time: time.Now(),
	}
	handler.storageWriter.storeLink(&linkInfo)

	if err := json.NewEncoder(w).Encode(&linkInfo); err != nil {
		log.Fatal("Error while parsing creating info return struct ", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		log.Print("Created In CreateLinkHandler with link: ", link)
	}
}

type DeleteLinkHandler struct {
	storageWriter *StorageWriter
}

func (handler DeleteLinkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var link_id string
	if linkList := r.URL.Query()["link_id"]; len(linkList) == 1 {
		link_id = string(linkList[0])
	} else {
		log.Print("Incoming bad request on delete link")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	found := handler.storageWriter.deleteLink(link_id)

	if found {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

type LinksInfoHandler struct {
	storageWriter *StorageWriter
}

func (handler LinksInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	links := handler.storageWriter.linksInfo()

	if err := json.NewEncoder(w).Encode(links); err != nil {
		log.Fatal("Error while encoding links info ", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		log.Print("Generated LinksInfo")
		w.Header().Set("Content-Type", "application/json")
	}
}