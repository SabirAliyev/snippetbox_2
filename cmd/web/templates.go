package main

import (
	"html/template"
	"path/filepath"
	"time"

	"sabiraliyev.net/snippetbox/pkg/forms"
	"sabiraliyev.net/snippetbox/pkg/models"
)

// Define a templateData type to act as the holding structure for any dynamic data we want to pass
// to our HTML templates.
type templateData struct {
	CSRFToken       string
	CurrentYear     int
	Flash           string
	Form            *forms.Form
	IsAuthenticated bool
	IsAdministrator bool
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Message         *models.Message
	Messages        []*models.Message
}

// Create a humanDate function which returns a nicely formatted string representation of a time.Time object.
func humanDate(t time.Time) string {
	// Return the empty string if time has the zero value.
	if t.IsZero() {
		return ""
	}

	// Convert the time to UTC before formatting it.
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in global variable. This is essentially a string-keyed
// map which acts as a lookup between the names of our custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with the extension '.page.tmpl'.
	// This essentially gives us a slice of all the 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop the pages one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl') from the full file path
		// and assign it to the variable.
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'layout' templates to the template set
		// (in our case, it`s just a 'base' layout at the moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'partial' templates to the template set
		// (in our case, it`s just a 'footer' partial at the moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page (like 'home.page.tmpl') as the key.
		cache[name] = ts
	}

	// Return the map
	return cache, nil
}
