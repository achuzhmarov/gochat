package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		if route.AuthHandlerFunc != nil {
			authHandler := route.AuthHandlerFunc
			handler = Authenticator(authHandler, route.Name)
		} else {
			handler = route.HandlerFunc
		}

		handler = Logger(handler, route.Name)
		handler = ErrorHandler(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
