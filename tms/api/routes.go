package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/tasks", app.createTaskHandler)
	router.HandlerFunc(http.MethodGet, "/tasks/:id", app.showTaskHandler)
	return router
}
