package main

import "sabiraliyev.net/snippetbox/pkg/models"

// Define a templateData type to act as the holding structure for any dynamic data we want to pass
// to our HTML templates. At the moment it only contains one field, but we'll add more to it
// in the build process.
type templateData struct {
	Snippet *models.Snippet
}
