package main;

import (
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

type JWTTokensMngr struct {
	hmacSecret []byte
}

func NewJWTTokensMngr(secret string) (mngr *JWTTokensMngr) {
	mngr = new(JWTTokensMngr)
	mngr.hmacSecret = []byte(secret)
	return
}

func (mngr JWTTokensMngr) createTokens(userInfo UserInfo) (acc string, refr string, err error) {
	created := time.Now().Unix()
	
	accToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"type": "access",
		"userID": userInfo.userID,
		"email": userInfo.email,
		"login": userInfo.login,
		"created": created,
	})
	if acc, err = accToken.SignedString(mngr.hmacSecret); err != nil {
		log.Fatal("Something went wrong while signing access token", err)
		return
	}

	refrToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
		"type": "refresh",
		"login": userInfo.login,
		"password": userInfo.password,
		"created": created,
	})
	if refr, err = refrToken.SignedString(mngr.hmacSecret); err != nil {
		log.Fatal("Something went wrong while signing refresh token", err)
	}
	return
}

type ValidationResult struct {
	result int
	userID int
	email string
	login string
}

const (
	ValidationResultOk int = 0
	ValidationResultInvalid int = 1
	ValidationResultExpired int = 2
)

const (
	AccessTokenAliveMs int64 = 1000 * 60 * 15		// 15min
	RefreshTokenAliveMs int64 = 1000 * 60 * 24 * 14	// 14 days
)

func (mngr JWTTokensMngr) validateAccessToken(acc string) (result ValidationResult) {
	token, err := jwt.Parse(acc, func(token *jwt.Token) (interface{}, error) {return mngr.hmacSecret, nil})
	if err != nil {
		log.Println("Something went wrong while parsing access token")
		result.result = ValidationResultInvalid
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok {
		log.Println("I cant get claims")
		result.result = ValidationResultInvalid
		return
	}

	if claims["type"] != "access" {
		log.Println("Not access token")
		result.result = ValidationResultInvalid
		return
	}

	created := int64(claims["created"].(float64))

	expired := created + AccessTokenAliveMs
	now := time.Now().Unix()
	if now < created || expired < now {
		log.Println("Got Expired token")
		result.result = ValidationResultExpired
		return
	}

	result.result = ValidationResultOk
	result.userID = int(claims["userID"].(float64))
	result.email = claims["email"].(string)
	result.login = claims["login"].(string)
	return
}

func (mngr JWTTokensMngr) parseRefreshToken(refrToken string) (login string, password string, refreshed bool) {
	refreshed = false
	token, err := jwt.Parse(refrToken, func(token *jwt.Token) (interface{}, error) {return mngr.hmacSecret, nil})
	if err != nil {
		log.Println("Something went wrong while parsing refresh token", err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok {
		log.Println("I cant get claims")
		return
	}

	if claims["type"] != "refresh" {
		log.Println("Not refresh token")
		return
	}

	created := int64(claims["created"].(float64))

	expired := created + RefreshTokenAliveMs
	now := time.Now().Unix()
	if now < created || expired < now {
		log.Println("Got Expired token")
		return
	}

	login = claims["login"].(string)
	password = claims["password"].(string)
	refreshed = true
	return
}