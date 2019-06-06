package main

import "cesarbon.net/goproject/pkg/models"

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
