package main

import (
	"flag"
	"log"
	"net/http"
)

func main(){

	//flag name = addr, default value = :4000, description of flag = HHTP...
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	//use the http.newservemux() function to initialize a new servemux,
	//then register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//static url for handles
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Starting server on ", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}