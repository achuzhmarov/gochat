package router

import (
	"gochat/web/controller"

	"net/http"

	log "github.com/Sirupsen/logrus"
)

func ErrorHandler(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(

		func(response http.ResponseWriter, request *http.Request) {
			defer computeError(response, request, name)

			inner.ServeHTTP(response, request)
		})
}

func computeError(response http.ResponseWriter, request *http.Request, name string) {
	err := recover()

	if err != nil {
		log.WithFields(log.Fields{
			"method":     request.Method,
			"requestURI": request.RequestURI,
			"name":       name,
			"err":        err,
		}).Error("Error in handling request")

		controller.AddPlainTextContentType(response)
		response.WriteHeader(http.StatusBadRequest)
	}
}
