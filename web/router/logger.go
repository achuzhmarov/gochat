package router

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(
		func(response http.ResponseWriter, request *http.Request) {
			start := time.Now()

			inner.ServeHTTP(response, request)

			log.WithFields(log.Fields{
				"method":     request.Method,
				"requestURI": request.RequestURI,
				"name":       name,
				"time":       time.Since(start),
			}).Info("handle request")
		},
	)
}
