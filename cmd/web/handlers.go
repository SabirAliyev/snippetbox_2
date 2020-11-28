package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Change the signature of the home handler so it is defined as a method against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	// Initialize a slice containing th e path to the two files. Notice that the
	// home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of the paths
	// as a variadic parameter.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		// Because the home handler function is now method against application
		// it can access its fields, including the error logger. We`ll write the log
		// message to this instead of the standard logger.
		// app.errorLog.Println(err.Error())
		// http.Error(w, "Internal Server Error", 500) // DEPRECATED

		app.serverError(w, err) // Use the serverError() helper.
		return
	}

	// Use Execute() method on the template set to write the template content as the response
	// body. The last parameter to Execute() represents any dynamic data that we want to pass in,
	// which for now we`ll leave as nil.
	err = ts.Execute(w, nil)
	if err != nil {
		// Also update the code here to use the error logger from the application struct.
		// app.errorLog.Println(err.Error()) // DEPRECATED

		app.serverError(w, err) // Use the serverError() helper.
		http.Error(w, "Internal Server Error", 500)
	}
}

// Change the signature of the showSnippet() handler so it is defined as a method against *application.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		// http.NotFound(w, r) // DEPRECATED

		app.notFound(w) // Use the notFound() helper.
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Change the signature of the createSnippet() handler so it is defined as a method against *application.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// http.Error(w, "Method not allowed", 405) // DEPRECATED

		app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
