package setlx2python_playground

import (
	"net/http"

	"github.com/gorilla/mux"
)

// CreateRouter maps RequestHandler functions to REST API
func CreateRouter(handler *RequestHandler) *mux.Router {

	router := mux.NewRouter()
	router.StrictSlash(true)

	// index page handler
	router.Path("/").Methods("GET").HandlerFunc(handler.index)

	router.Path("/transpile").Methods("POST").HandlerFunc(handler.transpile)

	router.Path("/run/setlx").Methods("POST").HandlerFunc(handler.runSetlX)
	// router.Path("/run/python").Methods("POST").HandlerFunc(handler.runPython)

	// serve static files
	fileHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("www/static")))
	router.PathPrefix("/static/").Handler(fileHandler)

	return router
}
