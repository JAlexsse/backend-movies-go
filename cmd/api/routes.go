package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)

	return app.enableCORS(router) //adds the middleware to all the requests
}
