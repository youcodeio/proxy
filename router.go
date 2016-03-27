package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

type Route struct {
Name        string
Method      string
Pattern     string
HandlerFunc http.HandlerFunc
}
type Routes []Route

var routes = Routes{
Route{"GetChannels", "GET", "/channels", Channelslist},
Route{"GetQuery", "GET", "/query", GetQuery},
}

func NewRouter() *mux.Router {
router := mux.NewRouter().StrictSlash(true)
for _, route := range routes {
router.
Methods(route.Method).Path(route.Pattern).
Name(route.Name).Handler(route.HandlerFunc)
}
return router
}
