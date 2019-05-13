package routers

import (
	"github.com/gorilla/mux"
	"net/http"
)

const STATIC_DIR = "/Users/victor/workspace/go/src/github.com/vikas91/consistent-file-access/application/static/"

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(STATIC_DIR))))
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}