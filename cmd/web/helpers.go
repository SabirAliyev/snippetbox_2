package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// The ServerError helper writes an error message and stack trace to the errorLog,
// then sends e generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	// To report the file name and line number *one step back* in the stack trace.
	// We do this by setting the frame depth to 2.
	app.errorLog.Println(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper send a specific status code and corresponding description
// to the user. We`ll use this later in the book to send responses like 400 "Bad Request"
// when there`s a problem with the request that the user send.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we`ll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which send a 404 Not Found response to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, _ *http.Request, name string, td *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name
	// (like 'home.page.tmpl'). If no entry exists in the cache with the provided name,
	// call the serverError helper that we made earlier.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the templaste %s does not exist", name))
		return
	}

	// Execute the template set, passing in any dynamic data.
	err := ts.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
	}
}
