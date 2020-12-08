package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sabiraliyev.net/snippetbox/pkg/models"
)

// Change the signature of the home handler so it is defined as a method against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because Pat matches the "/" path exactly, we can now remove the manual check of r.URL.PAth != "/" from this handler.

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the render helper.
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

// Change the signature of the showSnippet() handler so it is defined as a method against *application.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Pat doesn`t strip the colon from the names capture key,
	// so we need to get the value of ":id" from the query string instead of "id".
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	// Use the SnippetModel object`s Get method to retrieve the data for a specific record based on ID.
	// If no matching record is found, return a 404 Not Found response.
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Use the render helper.
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

// Add new createSnippetForm handler, which for now a placeholder response.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// First we call r.ParseForm() which add any data in POST request bodies to the r.PostForm map.
	// This also works in the same way for PUT and PATCH requests. If there any errors, we use
	// our app.Client helper to send a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the relevant data fields from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	// Create a new snippet record in the database using the form data.
	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
