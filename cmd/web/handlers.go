//this imports main
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"cesarbon.net/goproject/pkg/forms"
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

	// Use the PopString() method to retrieve the value for the "flash" key.
	// PopString() also deletes the key and value from the session data, so it
	// acts like a one-time fetch. If there is no matching key in the session
	// data this will return the empty string.
	//flash := app.session.PopString(r, "flash")

	// Use the new render helper.
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData {
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)

	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}	

	//now we can create the snippet model, since we know that input is valid.
	sniptemp := &models.Snip{}

	sniptemp.Title = r.PostForm.Get("title")
	sniptemp.Content = r.PostForm.Get("content")
	sniptemp.Expires = r.PostForm.Get("expires")
	//errors := make(map[string]string)


	// this logic should be set in its own class
	// and be called with an object e.g snippet.validate()
	// if strings.TrimSpace(sniptemp.Title) == "" {
	// 	errors["title"] = "this field cannot be blank"
	// } else if utf8.RuneCountInString(sniptemp.Title) > 100 {
	// 	errors["title"] = "This field is too long (maximum is 100 characters)"
	// }

	// if strings.TrimSpace(sniptemp.Content) == "" {
	// 	errors["content"] = "this field cannot be blank"
	// }

	// if strings.TrimSpace(sniptemp.Expires) == "" {
	// 	errors["expires"] = "this field cannot be blank"
	// } else if sniptemp.Expires != "365" && sniptemp.Expires != "7" && sniptemp.Expires != "1" {
	// 	errors["expires"] = "This field is invalid"
	// }

	// if len(errors) > 0 {
	// 	app.render(w, r, "create.page.tmpl", &templateData{
	// 		Form: forms.New(nil),
	// 	})
	// 	return
	// }

	// here we are sending the form's data to data layer for snippet object.
	id, err := app.snippets.Insert(*sniptemp)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the Put() method to add a string value ("Your snippet was saved
	// successfully!") and the corresponding key ("flash") to the session
	// data. Note that if there's no existing session for the current user
	// (or their session has expired) then a new, empty, session for them
	// will automatically be created by the session middleware.
	app.session.Put(r, "flash", "Snippet successfully created!")

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

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	//Parse form data
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))

	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData {
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//Check whether the credentials are valid. If they are not, add
	//a generic error message to the form failures map and re-display the login page
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "authenticatedUserID", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")

	app.session.Put(r, "flash", "You've been logged out successfully")
	http.Redirect(w, r, "/", 303)
}

