package controller

import (
	"gochat/logic"
	"gochat/model"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func RegisterNewUser(response http.ResponseWriter, request *http.Request) {
	var requestUser *model.UserAuth

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&requestUser)
	panicIfError(err)

	_, err = logic.UserService.RegisterNewUser(requestUser)

	_, alreadyExist := err.(*logic.UserAlreadyExists)

	if alreadyExist {
		AddJsonContentType(response)
		response.WriteHeader(http.StatusConflict)
		return
	}

	panicIfError(err)

	AddJsonContentType(response)
	response.WriteHeader(http.StatusOK)
}

func GetSelfInfo(response http.ResponseWriter, request *http.Request, userLogin string) {
	user, err := logic.UserService.GetUserByLogin(userLogin)

	if err != nil {
		panic(err)
	}

	AddJsonContentType(response)
	response.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(response).Encode(user); err != nil {
		panic(err)
	}
}

func GetUnreadTalks(response http.ResponseWriter, request *http.Request, userLogin string) {
	talks, err := logic.UserService.GetUnreadTalks(userLogin)
	panicIfError(err)

	AddJsonContentType(response)
	response.WriteHeader(http.StatusOK)

	if talks != nil && len(talks) != 0 {
		if err := json.NewEncoder(response).Encode(talks); err != nil {
			panic(err)
		}
	}
}

type userSuggestions struct {
	Suggestions int `json:"suggestions"`
}

func UseSuggestions(response http.ResponseWriter, request *http.Request, userLogin string) {
	vars := mux.Vars(request)
	countString := vars["count"]

	count, err := strconv.Atoi(countString)
	panicIfError(err)

	_, err = logic.UserService.UseSuggestions(userLogin, count)
	panicIfError(err)

	AddJsonContentType(response)
	response.WriteHeader(http.StatusOK)
}

func GetNewFriends(response http.ResponseWriter, request *http.Request, userLogin string) {
	vars := mux.Vars(request)
	searchString := vars["searchString"]

	friends, err := logic.UserService.GetNewFriends(userLogin, searchString)
	panicIfError(err)

	AddJsonContentType(response)
	response.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(response).Encode(friends); err != nil {
		panic(err)
	}
}

func MakeFriends(response http.ResponseWriter, request *http.Request, userLogin string) {
	vars := mux.Vars(request)
	friendLogin := vars["friendLogin"]

	talk, err := logic.UserService.MakeFriends(userLogin, friendLogin)
	panicIfError(err)

	AddJsonContentType(response)
	response.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(response).Encode(talk); err != nil {
		panic(err)
	}
}
