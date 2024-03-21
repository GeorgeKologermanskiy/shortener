package main

import (
	"encoding/json"
	"fmt"
    "net/http"
    "net/http/httptest"
    //"strings"
    "testing"
	"time"

    "github.com/gorilla/mux"
)

func run(r *mux.Router, req *http.Request, t *testing.T) (wr *httptest.ResponseRecorder) {
	wr = httptest.NewRecorder()
	var match mux.RouteMatch
	if !r.Match(req, &match) {
		t.Error("Route was not found?!?")
	}
	match.Handler.ServeHTTP(wr, req)
	return
}

func TestCreateRequest(t *testing.T) {
	test_host := "test_host"
	storageWriter := newStorageWriter(test_host)
	r := createRouter(storageWriter)
	
	clientLink := "12345"
	clientURL := fmt.Sprintf("/create_link?link=%s", clientLink)
	req := httptest.NewRequest(http.MethodPost, clientURL, nil)
	timeBeforeRun := time.Now()
	httpResp := run(r, req, t)

	if httpResp.Code != http.StatusOK {
		t.Error("expected HTTP status 200, got ", httpResp.Code)
		return
	}
	var resp ShortLinkInfo
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		t.Error("Failed while parsing response ", err)
		return
	}
	if len(resp.Link_id) != 7 {
		t.Error("Invalid Link_id len")
		return
	}
	if resp.Short_link != test_host + "/" + resp.Link_id {
		t.Errorf("Invalid Short_link, expected %s, got %s", test_host + "/" + resp.Link_id, resp.Short_link)
		return
	}
	if resp.Full_link != clientLink {
		t.Errorf("Invalid Full_link, expected %s, got %s", test_host + "/" + clientLink, resp.Full_link)
		return
	}
	if resp.Created_time.Before(timeBeforeRun) {
		t.Error("HOW??!?!??")
		return
	}
	if resp.Expired_time.Before(resp.Expired_time) {
		t.Error("HOW[2]??!?!??")
		return
	}
}

func TestCreateAndDeleteRequest(t *testing.T) {
	test_host := "test_host"
	storageWriter := newStorageWriter(test_host)
	r := createRouter(storageWriter)
	
	clientLink := "12345"
	clientURL := fmt.Sprintf("/create_link?link=%s", clientLink)
	req := httptest.NewRequest(http.MethodPost, clientURL, nil)
	httpResp := run(r, req, t)

	var resp ShortLinkInfo
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		t.Error("Failed while parsing response ", err)
		return
	}
	link_id := resp.Link_id

	clientURL = fmt.Sprintf("/delete_link?link_id=%s", link_id)
	req = httptest.NewRequest(http.MethodDelete, clientURL, nil)
	httpResp = run(r, req, t)

	if httpResp.Code != http.StatusOK {
		t.Error("expected HTTP status 200, got ", httpResp.Code)
	}

}