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

type ShortLinkInfo struct {
	Link_id string `json:"link_id"`
	Short_link string `json:"short_link"`
	Full_link string `json:"full_link"`
	Created_time time.Time `json:"created_time"`
	Expired_time time.Time `json:"expired_time"`
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
	link_id, short_link := handler.storageWriter.addLink(link)

	if err := json.NewEncoder(w).Encode(&ShortLinkInfo {
		Link_id: link_id,
		Short_link: short_link,
		Full_link: link,
		Created_time: time.Now(),
		Expired_time: time.Now(),
	}); err != nil {
		log.Fatal("Error while parsing return struct ", err)
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