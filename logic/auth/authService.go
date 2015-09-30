package auth

import (
	"gochat/model"

	"encoding/json"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

type TokenAuthentication struct {
	Token string `json:"token" form:"token"`
}

func Login(requestUser *model.UserAuth) (int, []byte) {
	authBackend := InitJWTAuthenticationBackend()

	if authBackend.Authenticate(requestUser) {
		token, err := authBackend.GenerateToken(requestUser.Login)

		if err != nil {
			return http.StatusInternalServerError, []byte("")
		} else {
			response, _ := json.Marshal(TokenAuthentication{token})
			return http.StatusOK, response
		}
	}

	return http.StatusUnauthorized, []byte("")
}

func RefreshToken(userLogin string) []byte {
	authBackend := InitJWTAuthenticationBackend()

	token, err := authBackend.GenerateToken(userLogin)
	if err != nil {
		panic(err)
	}

	response, err := json.Marshal(TokenAuthentication{token})
	if err != nil {
		panic(err)
	}

	return response
}

func Logout(req *http.Request) error {
	authBackend := InitJWTAuthenticationBackend()

	tokenRequest, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	if err != nil {
		return err
	}

	tokenString := req.Header.Get("Authorization")

	return authBackend.Logout(tokenString, tokenRequest)
}
