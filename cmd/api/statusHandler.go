package main

import (
	"net/http"
)

func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := AppStatus{
		Status:     "Available",
		Enviroment: app.config.env,
		Version:    version,
	}

	app.writeJSON(w, http.StatusOK, currentStatus, "status")
}
