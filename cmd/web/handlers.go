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
	// Because Pat matches the "/" path exactly, ewe can now remove the manual check of r.URL.PAth != "/" from this handler.

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
	// so we need to get the value of ":d" from the query string instead of "id".
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
func (app *application) createSnippetForm(w http.ResponseWriter, _ http.Request) {
	w.Write([]byte("Create a new snippet..."))
}

// Change the signature of the createSnippet() handler so it is defined as a method against *application.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Checking if the request method is a POST is now superfluous and can be removed.

	// Create some variables holding dummy data.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := "7"

	// Pass the data to the SnippetModel.Insert() method, receiving the ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Change the redirect to use the new semantic URL style of /snippet:/:id.
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
