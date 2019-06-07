//this imports main
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"cesarbon.net/goproject/cmd/config"
	"cesarbon.net/goproject/pkg/models"
)

//Home function from handlers
func Home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// files := []string{
		// 	"./ui/html/home.page.tmpl",
		// 	"./ui/html/base.layout.tmpl",
		// 	"./ui/html/footer.partial.tmpl",
		// }

		// ts, err := template.ParseFiles(files...)
		// if err != nil{
		// 	app.ErrorLog.Println(err.Error())
		// 	http.Error(w, "Internal Server Error_1", 500)
		// 	return
		// }

		// err = ts.Execute(w, nil)
		// if err != nil{
		// 	app.ErrorLog.Println(err.Error())
		// 	http.Error(w, "Internal Server Error_2", 500)
		// }
	}
}

//define a home handle function which writes a byte slice containing
//"hello from snippetbox" as the response.
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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

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

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

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

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
