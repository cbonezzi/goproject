package main

import (
	"html/template"
	"path/filepath"
	"time"

	"cesarbon.net/goproject/pkg/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

//this is a custom template function
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

//function that get store in a global variable. essentially a
//string-keyed map which acts as a lookup between the names
//of our custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {

	//init a new map to act as the cache
	cache := map[string]*template.Template{}

	//get a slice of all filepaths with the extension '.page.tmpl'
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	//loop through the pages
	for _, page := range pages {
		name := filepath.Base(page)

		//the template.FuncMap must be registered with the template
		//before we call the parsefiles() method
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
