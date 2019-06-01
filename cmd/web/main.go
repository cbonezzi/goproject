package main

import (	
	"cesarbon.net/goproject/cmd/config"
	"flag"
	"net/http"
	"os"
	"log"
)

//Cfg struct
type Cfg struct{
	Addr string
	StaticDir string
	IsProd bool
}

//struct for dependency injection
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main(){

	cfg := new(Cfg)
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address1")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static","Path to static assests")
	flag.BoolVar(&cfg.IsProd, "isProd", true, "flag for production enable")

	flag.Parse()

	//way of creating log distinction and separation...perhaps put this into a class,
	//and abstract it out, to use it globally?
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	
	//di but basic
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	//init var with struct reqs for di.
	app1 := &config.Application{
		InfoLog: log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	//use the http.newservemux() function to initialize a new servemux,
	//then register the home function as the handler for the "/" URL pattern.
	//DI: switching home -> app.home since method now support DI (basic)
	mux := http.NewServeMux()
	mux.Handle("/", Home(app1))
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//static url for handles
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	
	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr: cfg.Addr,
		ErrorLog: app1.ErrorLog,
		Handler: mux,
	}

	app1.InfoLog.Printf("Starting server on %s", cfg.Addr)
	app1.InfoLog.Println("static directory location", cfg.StaticDir)
	app1.InfoLog.Println("Is Production", cfg.IsProd)
	err := srv.ListenAndServe()
	app1.ErrorLog.Fatal(err)
}