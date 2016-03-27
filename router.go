package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route struct
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes slice of routes
type Routes []Route

var routes = Routes{
	Route{"GetChannels", "GET", "/channels", Channelslist},
	Route{"GetQuery", "GET", "/query", GetQuery},
}

//NewRouter return a instance of mux.Router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).Path(route.Pattern).
			Name(route.Name).Handler(route.HandlerFunc)
	}
	return router
}
