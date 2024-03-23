package main

import (
	"net/http"
    "net/http/httptest"
    "testing"
	
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

func TestRegisterLoginValidateRefreshRequest(t *testing.T) {
	us := NewUserStorage()
	jwtMngr := NewJWTTokensMngr("123")
	r := createRouter(us, jwtMngr)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", nil)
	req.Header.Add("X-USER-EMAIL", "gg@gg.gg")
	req.Header.Add("X-USER-LOGIN", "gg1")
	req.Header.Add("X-USER-PASSWORD", "gg2")
	httpResp := run(r, req, t)
	if httpResp.Code != http.StatusOK {
		t.Error("expected HTTP status 200, got ", httpResp.Code)
		return
	}

	req = httptest.NewRequest(http.MethodGet, "/api/auth/auth", nil)
	req.Header.Add("X-USER-LOGIN", "gg1")
	req.Header.Add("X-USER-PASSWORD", "gg2")
	httpResp = run(r, req, t)
	if httpResp.Code != http.StatusOK {
		t.Error("expected HTTP status 200, got ", httpResp.Code)
		return
	}
	acc := httpResp.Header().Get("X-ACCESS-TOKEN")
	refr := httpResp.Header().Get("X-REFRESH-TOKEN")

	req = httptest.NewRequest(http.MethodGet, "/api/auth/validate", nil)
	req.Header.Add("X-ACCESS-TOKEN", acc)
	httpResp = run(r, req, t)
	if httpResp.Code != http.StatusOK {
		t.Error("expected HTTP status 200, got ", httpResp.Code)
		return
	}

	req = httptest.NewRequest(http.MethodGet, "/api/auth/validate", nil)
	req.Header.Add("X-ACCESS-TOKEN", refr)
	httpResp = run(r, req, t)
	if httpResp.Code != http.StatusUnauthorized {
		t.Error("expected HTTP status 401, got ", httpResp.Code)
		return
	}

	req = httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
	req.Header.Add("X-REFRESH-TOKEN", refr)
	httpResp = run(r, req, t)
	if httpResp.Code != http.StatusOK {
		t.Error("expected HTTP status 200, got ", httpResp.Code)
		return
	}

	req = httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
	req.Header.Add("X-REFRESH-TOKEN", acc)
	httpResp = run(r, req, t)
	if httpResp.Code != http.StatusUnauthorized {
		t.Error("expected HTTP status 401, got ", httpResp.Code)
		return
	}
}
