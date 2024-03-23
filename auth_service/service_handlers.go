package main;

import (
	"log"
	"net/http"
    "strconv"
)

type RegisterHandler struct {
	us *UserStorage
}

func (handler RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("X-USER-EMAIL")
	login := r.Header.Get("X-USER-LOGIN")
	password := r.Header.Get("X-USER-PASSWORD")
	// TODO: Validate email, login, password

	exists := handler.us.addUser(email, login, password)
	if exists {
		log.Println("Found user with same login:", login)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("Registered User", email, login)
}

type AuthentificateHandler struct {
	us *UserStorage
	jwtMngr *JWTTokensMngr
}

func (handler AuthentificateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	login := r.Header.Get("X-USER-LOGIN")
	password := r.Header.Get("X-USER-PASSWORD")

	info, found := handler.us.getUser(login)
	if !found {
		log.Println("User", login, "not found")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if info.password != password {
		log.Println("Got invalid password for", login)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	
	log.Println("Successfully found user", login)
	if acc, refr, err := handler.jwtMngr.createTokens(info); err != nil {
		log.Fatalln("Error while creating tokens", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		log.Println("Correctly created tokens")
		w.Header().Add("X-ACCESS-TOKEN", acc)
		w.Header().Add("X-REFRESH-TOKEN", refr)	
	}
}

type ValidateHandler struct {
	jwtMngr *JWTTokensMngr
}

func (handler ValidateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	accToken := r.Header.Get("X-ACCESS-TOKEN")
	validateResult := handler.jwtMngr.validateAccessToken(accToken)

	if validateResult.result != ValidationResultOk {
		w.Header().Add("AUTH-RESULT-INFO", strconv.Itoa(validateResult.result))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	
	log.Println("Token validated, OK")
	w.Header().Add("X-USER-USERID", strconv.Itoa(validateResult.userID))
	w.Header().Add("X-USER-EMAIL", validateResult.email)
	w.Header().Add("X-USER-LOGIN", validateResult.login)
	w.WriteHeader(http.StatusOK)
}

type RefreshHandler struct {
	us *UserStorage
	jwtMngr *JWTTokensMngr
}

func (handler RefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	refrToken := r.Header.Get("X-REFRESH-TOKEN")
	login, password, refreshed := handler.jwtMngr.parseRefreshToken(refrToken)
	if !refreshed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userInfo, ok := handler.us.getUser(login)
	if !ok || userInfo.password != password {
		log.Println("Invalid user in refresh token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if acc, refr, err := handler.jwtMngr.createTokens(userInfo); err != nil {
		log.Println("Error while creating new tokens", err)	
		w.WriteHeader(http.StatusUnauthorized)	
	} else {
		log.Println("Token refreshed, OK")
		w.Header().Add("X-ACCESS-TOKEN", acc)
		w.Header().Add("X-REFRESH-TOKEN", refr)
		w.WriteHeader(http.StatusOK)
	}
}
