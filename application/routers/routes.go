package routers

import (
	"github.com/vikas91/consistent-file-access/application/handlers"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		handlers.Index,
	},
	Route{
		"RegisterNode",
		"POST",
		"/register",
		handlers.RegisterNode,
	},
	Route{
		"ShowNodeList",
		"GET",
		"/nodes",
		handlers.ShowNodeList,
	},
}