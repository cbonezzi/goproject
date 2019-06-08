package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//use the http.newservemux() function to initialize a new servemux,
	//then register the home function as the handler for the "/" URL pattern.
	//DI: switching home -> app.home since method now support DI (basic)
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("F:/snippetbox_git/goproject/ui/static/"))

	//static url for handles
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Pass the servemux as the 'next' parameter to the secureHeaders middleware.
	// Because secureHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.

	//wrap the existing chain with the logRequest middleware

	// wrap the exisiting chain with app.recoverPanic to handle the panic gracefully.
	//without alice
	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	//with alice, return the 'standard' middleware
	return standardMiddleware.Then(mux)
}
