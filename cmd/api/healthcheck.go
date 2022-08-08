package main

import (
	"net/http"
)

// Encoding of nested objects is supported. So, for example, if you have a slice of structs in
// Go that will encode to an array of objects in JSON.
// `
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	env := envelope{
		"status": "available", "system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, env)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
