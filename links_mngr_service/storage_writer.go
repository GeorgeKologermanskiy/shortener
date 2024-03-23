package main

import (
	"math/rand"
	"time"
)

type ShortLinkInfo struct {
	Link_id string `json:"link_id"`
	Short_link string `json:"short_link"`
	Full_link string `json:"full_link"`
	Created_time time.Time `json:"created_time"`
	Expired_time time.Time `json:"expired_time"`
}

type StorageWriter struct {
	links map[string]ShortLinkInfo
	proxy_host string
}

func newStorageWriter(proxy_host string) (storageWriter *StorageWriter) {
	storageWriter = new(StorageWriter)
	storageWriter.links = make(map[string]ShortLinkInfo)
	storageWriter.proxy_host = proxy_host
	return
}

var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func generateLinkId() string {
	res := make([]byte, 7)
	for i := range res {
		res[i] = charset[rand.Intn(len(charset))]
	}
	return string(res)
}

func (s StorageWriter) storeLink(linkInfo *ShortLinkInfo) {
	for {
		link_id := generateLinkId()
		if _, ok := s.links[link_id]; ok {
			continue
		}

		linkInfo.Link_id = link_id
		linkInfo.Short_link = s.proxy_host + "/" + link_id
		s.links[link_id] = *linkInfo
		return
	}
}

func (s StorageWriter) deleteLink(link_id string) bool {
	if _, ok := s.links[link_id]; ok {
		delete(s.links, link_id)
		return true
	}

	return false
}

func (s StorageWriter) linksInfo() []ShortLinkInfo {
	res := make([]ShortLinkInfo, 0)
	for _, v := range s.links {
		res = append(res, v)
	}
	return res
}
