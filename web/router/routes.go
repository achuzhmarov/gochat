package router

import (
	"gochat/web/controller"

	"net/http"
)

type Route struct {
	Name            string
	Method          string
	Pattern         string
	AuthHandlerFunc HandlerWithAuth
	HandlerFunc     http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	createRoute("TokenAuth", "POST", "/login", controller.Login),
	createRoute("Registration", "POST", "/user/add", controller.RegisterNewUser),
	createAuthRoute("RefreshToken", "POST", "/login_refresh", controller.RefreshToken),
	createAuthRoute("TokenAuth", "POST", "/logout", controller.Logout),

	createAuthRoute("GetSelfInfo", "GET", "/user", controller.GetSelfInfo),
	createAuthRoute("GetUnreadTalks", "GET", "/user/new_talks", controller.GetUnreadTalks),
	createAuthRoute("UseSuggestions", "POST", "/user/use_suggestions/{count}", controller.UseSuggestions),
	createAuthRoute("GetNewFriends", "GET", "/user/friends/{searchString}", controller.GetNewFriends),
	createAuthRoute("GetRandomNewFriends", "GET", "/user/friends", controller.GetNewFriends),
	createAuthRoute("MakeFriends", "POST", "/user/friend/{friendLogin}", controller.MakeFriends),

	createAuthRoute("GetMessagesByTalk", "GET", "/messages/{talkId}", controller.GetMessagesByTalk),
	createAuthRoute("GetUnreadMessagesByTalk", "GET", "/messages/new/{talkId}", controller.GetUnreadMessagesByTalk),
	createAuthRoute("CreateMessage", "POST", "/messages/add", controller.CreateMessage),
	createAuthRoute("DecipherMessage", "POST", "/messages/decipher", controller.DecipherMessage),
}

func createRoute(name string, method string, pattern string, handler http.HandlerFunc) Route {
	return Route{name, method, pattern, nil, handler}
}

func createAuthRoute(name string, method string, pattern string, handler HandlerWithAuth) Route {
	return Route{name, method, pattern, handler, nil}
}
