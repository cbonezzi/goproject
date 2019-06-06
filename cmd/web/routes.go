package main

import "net/http"

func (app *application) routes() *http.ServeMux {

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

	return mux
}
