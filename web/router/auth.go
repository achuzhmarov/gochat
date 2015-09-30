package router

import (
	"gochat/logic/auth"

	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
)

type HandlerWithAuth func(http.ResponseWriter, *http.Request, string)

func (f HandlerWithAuth) ServeHTTP(response http.ResponseWriter, request *http.Request, userLogin string) {
	f(response, request, userLogin)
}

func Authenticator(inner HandlerWithAuth, name string) http.Handler {
	return http.HandlerFunc(
		func(response http.ResponseWriter, request *http.Request) {
			authBackend := auth.InitJWTAuthenticationBackend()

			token, err := jwt.ParseFromRequest(request, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				} else {
					return authBackend.PublicKey, nil
				}
			})

			isInBlackList := authBackend.IsInBlacklist(request.Header.Get("Authorization"))

			if err == nil && token.Valid && !isInBlackList {
				inner.ServeHTTP(response, request, token.Claims["sub"].(string))
			} else {
				log.WithFields(log.Fields{
					"err":           err.Error(),
					"token":         token,
					"isInBlackList": isInBlackList,
					"name":          name,
				}).Error("Unauthorized access")

				response.WriteHeader(http.StatusUnauthorized)
			}
		})
}
