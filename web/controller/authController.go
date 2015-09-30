package controller

import (
	"gochat/logic/auth"
	"gochat/model"

	"encoding/json"
	"net/http"
)

func Login(response http.ResponseWriter, request *http.Request) {
	var requestUser *model.UserAuth

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&requestUser)
	panicIfError(err)

	responseStatus, token := auth.Login(requestUser)
	AddJsonContentType(response)
	response.WriteHeader(responseStatus)
	response.Write(token)
}

func RefreshToken(response http.ResponseWriter, request *http.Request, userLogin string) {
	response.Header().Set("Content-Type", "application/json")
	response.Write(auth.RefreshToken(userLogin))
}

func Logout(response http.ResponseWriter, request *http.Request, userLogin string) {
	err := auth.Logout(request)
	panicIfError(err)

	AddPlainTextContentType(response)
	response.WriteHeader(http.StatusOK)
}
