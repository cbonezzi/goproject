package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

//Config struct
type Config struct{
	Addr string
	StaticDir string
}

func main(){

	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address1")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static","Path to static assests")

	//flag name = addr, default value = :4000, description of flag = HHTP...
	//addr := flag.String("addr", ":4000", "HTTP network address")
	//isProd := flag.Bool("IsProd", true, "Flag to determine if prod")

	flag.Parse()

	//way of creating log distinction and separation...perhaps put this into a class,
	//and abstract it out, to use it globally?
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//use the http.newservemux() function to initialize a new servemux,
	//then register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//static url for handles
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	infoLog.Printf("Starting server on %s", cfg.Addr)
	log.Println("Is Production", &cfg.StaticDir)
	err := http.ListenAndServe(cfg.Addr, mux)
	errorLog.Fatal(err)
}