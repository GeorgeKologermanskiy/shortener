package main

import (
	"math/rand"
)

type StorageWriter struct {
	links map[string]string
	proxy_host string
}

func newStorageWriter(proxy_host string) *StorageWriter {
	storageWriter := new(StorageWriter)
	storageWriter.links = make(map[string]string)
	storageWriter.proxy_host = proxy_host
	return storageWriter
}

var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func generateLinkId() string {
	res := make([]byte, 7)
	for i := range res {
		res[i] = charset[rand.Intn(len(charset))]
	}
	return string(res)
}

func (s StorageWriter) addLink(link string) (string, string) {
	for {
		link_id := generateLinkId()

		if _, ok := s.links[link_id]; ok {
			continue
		}

		s.links[link_id] = link
		return link_id, s.proxy_host + "/" + link_id
	}
}

func (s StorageWriter) getLink(link_id string) string {
	if link, ok := s.links[link_id]; ok {
		return link
	}

	return ""
}

func (s StorageWriter) deleteLink(link_id string) bool {
	if _, ok := s.links[link_id]; ok {
		delete(s.links, link_id)
		return true
	}

	return false
}
