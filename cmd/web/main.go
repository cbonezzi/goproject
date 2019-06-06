package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"cesarbon.net/goproject/cmd/config"

	"cesarbon.net/goproject/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

//Cfg struct
type Cfg struct {
	Addr      string
	StaticDir string
	IsProd    bool
	Dsn       string
}

//struct for dependency injection
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
}

func main() {

	cfg := new(Cfg)
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address1")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assests")
	flag.BoolVar(&cfg.IsProd, "isProd", true, "flag for production enable")
	flag.StringVar(&cfg.Dsn, "dsn", "web:P@ssw0rd@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	//way of creating log distinction and separation...perhaps put this into a class,
	//and abstract it out, to use it globally?
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//db connection pool creation.
	db, err := openDB(cfg.Dsn)

	//di but basic
	//initialize a mysql.SnippetModel instance and add it to the application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db},
	}

	if err != nil {
		app.errorLog.Fatal(err)
	}
	//defer a call to db.Close() so that the connection pool is closed before the main() function exits.
	defer db.Close()

	//init var with struct reqs for di.
	app1 := &config.Application{
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: app1.ErrorLog,
		Handler:  app.routes(),
	}

	app1.InfoLog.Printf("Starting server on %s", cfg.Addr)
	app1.InfoLog.Println("static directory location", cfg.StaticDir)
	app1.InfoLog.Println("Is Production", cfg.IsProd)
	err1 := srv.ListenAndServe()
	app1.ErrorLog.Fatal(err1)
}

func openDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
