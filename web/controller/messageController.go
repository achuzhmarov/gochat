package controller

import (
	"gochat/logic"
	"gochat/model"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetMessagesByTalk(response http.ResponseWriter, request *http.Request, userLogin string) {
	vars := mux.Vars(request)
	talkId := vars["talkId"]

	messages, err := logic.MessageService.GetMessagesByTalk(userLogin, talkId)
	panicIfError(err)

	AddJsonContentType(response)
	response.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(response).Encode(messages); err != nil {
		panic(err)
	}
}

func GetUnreadMessagesByTalk(response http.ResponseWriter, request *http.Request, userLogin string) {
	vars := mux.Vars(request)
	talkId := vars["talkId"]

	messages, err := logic.MessageService.GetUnreadMessagesByTalk(userLogin, talkId)
	panicIfError(err)

	AddJsonContentType(response)
	response.WriteHeader(http.StatusOK)

	if messages != nil && len(*messages) != 0 {
		if err := json.NewEncoder(response).Encode(messages); err != nil {
			panic(err)
		}
	}
}

func CreateMessage(response http.ResponseWriter, request *http.Request, userLogin string) {
	var message *model.Message

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&message)
	panicIfError(err)

	savedMessage, err := logic.MessageService.CreateMessage(userLogin, message)
	panicIfError(err)

	AddJsonContentType(response)
	response.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(response).Encode(savedMessage); err != nil {
		panic(err)
	}
}

func DecipherMessage(response http.ResponseWriter, request *http.Request, userLogin string) {
	var message *model.Message

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&message)
	panicIfError(err)

	_, err = logic.MessageService.DecipherMessage(userLogin, message)
	panicIfError(err)

	AddJsonContentType(response)
	response.WriteHeader(http.StatusOK)
}
