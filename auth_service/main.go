package main

import (
	"go.deanishe.net/env"
	"log"
	"net/http"

    "github.com/gorilla/mux"
)

func createRouter(us *UserStorage, jwtMngr *JWTTokensMngr) (r *mux.Router) {
	r = mux.NewRouter()
	s := r.PathPrefix("/api/auth").Subrouter()
	s.Handle("/register", RegisterHandler { us }).Methods(http.MethodPost)
	s.Handle("/auth", AuthentificateHandler { us, jwtMngr }).Methods(http.MethodGet)
	s.Handle("/validate", ValidateHandler { jwtMngr }).Methods(http.MethodGet)
	s.Handle("/refresh", RefreshHandler { us, jwtMngr }).Methods(http.MethodPost)
	return
}

func main() {
	log.Print("Starting service...")

	log.Print("Creating users storage...")
	us := NewUserStorage()

	log.Print("Creating JWTTokensMngr...")
	secret := env.Get("JWT_SECRET_KEY", "123456")
	if len(secret) == 0 {
		log.Fatal("Not found var JWT_SECRET_KEY")
		return
	}
	jwtMngr := NewJWTTokensMngr(secret)

	log.Print("Creating router with handlers...")
	r := createRouter(us, jwtMngr)

	log.Print("Start serving...")
	if err := http.ListenAndServe(":9000", r); err != nil {
		log.Fatal(err)
	}
}