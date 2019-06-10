//this imports main
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"cesarbon.net/goproject/pkg/models"
)

//define a home handle function
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// this helps extract the value of the id parameter from the query string
	// and try to convert it to an integer using the strconv.Atoi()
	// function. If it can't be converted to an integer, or the value is less than 1,
	// we return a 404 page not found response.
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// capturing Snippet struct
	s, err := app.snippets.Get(id)
	if err != nil {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
	}

	// Use the new render helper.
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	sniptemp := &models.Snip{}

	sniptemp.Title = r.PostForm.Get("title")
	sniptemp.Content = r.PostForm.Get("content")
	sniptemp.Expires = r.PostForm.Get("expires")

	errors := make(map[string]string)

	if strings.TrimSpace(sniptemp.Title) == "" {
		errors["title"] = "this field cannot be blank"
	} else if utf8.RuneCountInString(sniptemp.Title) > 100 {
		errors["title"] = "This field is too long (maximum is 100 characters)"
	}

	if strings.TrimSpace(sniptemp.Content) == "" {
		errors["content"] = "this field cannot be blank"
	}

	if strings.TrimSpace(sniptemp.Expires) == "" {
		errors["expires"] = "this field cannot be blank"
	} else if sniptemp.Expires != "365" && sniptemp.Expires != "7" && sniptemp.Expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	if len(errors) > 0 {
		fmt.Fprint(w, errors)
		return
	}

	id, err := app.snippets.Insert(*sniptemp)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetJSON(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		//header map manipulation occurs with the below --w.Header().Set()
		//w.Header().Add("Cache-Control", "public")
		//suppressing system-generated headers
		//w.Header()["Date"] = nil
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		w.WriteHeader(405)
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	//creates an instance of the snip object model
	sniptemp := &models.Snip{}

	//reads all bytes from request body and stores it into a slice
	jsn, err := ioutil.ReadAll(r.Body)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//unmarshal the jsn slice into the snip object model 'sniptemp'
	errUnmarshal := json.Unmarshal(jsn, sniptemp)
	if errUnmarshal != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//use sniptemp data from prop...
	id, err := app.snippets.Insert(*sniptemp)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
