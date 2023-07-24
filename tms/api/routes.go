package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/tasks", app.createTaskHandler2)
	router.HandlerFunc(http.MethodGet, "/tasks", app.getAllTasksHandler)
	router.Handle(http.MethodGet, "/tasks/:id", httprouter.Handle(app.getTaskHandler))
	router.Handle(http.MethodPut, "/tasks/:id", httprouter.Handle(app.updateTaskHandler))
	router.Handle(http.MethodDelete, "/tasks/:id", httprouter.Handle(app.deleteTaskHandler))

	return router
}
