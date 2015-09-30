package main

import (
	dal "gochat/dal/base"
	"gochat/web/router"

	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/weekface/mgorus"
)

func main() {
	hooker, err := mgorus.NewHooker(dal.DB_URL, dal.DB_NAME, "logs")
	if err == nil {
		log.AddHook(hooker)
	}

	webRouter := router.NewRouter()

	log.Info("Starting server")
	log.Fatal(http.ListenAndServe(":5000", webRouter))
}
