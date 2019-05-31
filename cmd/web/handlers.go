package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

//define a home handle function which writes a byte slice containing
//"hello from snippetbox" as the response.

func home(w http.ResponseWriter, r *http.Request){

	if r.URL.Path != "/"{
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil{
		log.Println(err.Error())
		http.Error(w, "Internal Server Error_1", 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil{
		log.Println(err.Error())
		http.Error(w, "Internal Server Error_2", 500)
	}
}

func showSnippet(w http.ResponseWriter, r *http.Request){
	// this helps extract the value of the id parameter from the query string
	// and try to convert it to an integer using the strconv.Atoi()
	// function. If it can't be converted to an integer, or the value is less than 1,
	// we return a 404 page not found response.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1{
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func createSnippet(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST"{
		//header map manipulation occurs with the below --w.Header().Set()
		//w.Header().Add("Cache-Control", "public")
		//suppressing system-generated headers
		//w.Header()["Date"] = nil
		w.Header().Set("Allow", "POST")
		w.WriteHeader(405)
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	w.Write([]byte("Create a new snippet..."))
}