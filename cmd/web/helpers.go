package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
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

// Create an addDefaultData helper. This takes a pointer to a TemplateData struct, add the current year
// to the CurrentYear field, and then returns the pointer. Again, we`re not using the *http.Request
// parameter at the moment.
func (app *application) addDefaultData(td *templateData, _ *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name
	// (like 'home.page.tmpl'). If no entry exists in the cache with the provided name,
	// call the serverError helper that we made earlier.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the templaste %s does not exist", name))
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Write the template to the buffer (buf), instead of straight to the http.ResponseWriter.
	// If there is an error, call our serverError helper and then return.
	// Execute the template set, passing the dynamic data with current year injected (app.addDefaultData())
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write the contents of the buffer to the http.ResponseWriter. Again, this is another time
	// where we pass our http.ResponseWriter to a function that takes an io.Writer.
	buf.WriteTo(w)
}