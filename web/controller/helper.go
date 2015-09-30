package controller

import (
	"net/http"
)

func AddJsonContentType(response http.ResponseWriter) {
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func AddPlainTextContentType(response http.ResponseWriter) {
	response.Header().Set("Content-Type", "text/plain; charset=UTF-8")
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
