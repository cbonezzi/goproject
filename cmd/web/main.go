package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"cesarbon.net/goproject/cmd/config"
	"cesarbon.net/goproject/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
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
	errorLog      *log.Logger
	infoLog       *log.Logger
	session		  *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	cfg := new(Cfg)
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address1")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assests")
	flag.BoolVar(&cfg.IsProd, "isProd", true, "flag for production enable")
	flag.StringVar(&cfg.Dsn, "dsn", "web:P@ssw0rd@/snippetbox?parseTime=true", "MySQL data source name")
	
	// Define a new command-line flag for the session secret (a random key which
	// will be used to encrypt and authenticate session cookies). It should be 32
	// bytes long.
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")

	flag.Parse()

	//way of creating log distinction and separation...perhaps put this into a class,
	//and abstract it out, to use it globally?
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//init template cache
	templateCache, err := newTemplateCache("F:/snippetbox_git/goproject/ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Use the sessions.New() function to initialize a new session manager,
	// passing in the secret key as the parameter. Then we configure it so
	// sessions always expires after 12 hours.
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	//db connection pool creation.
	db, err := openDB(cfg.Dsn)

	//di but basic
	//initialize a mysql.SnippetModel instance and add it to the application
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:	   session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
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

	//with this block of code we are telling Go that we should use Go's favored
	//cipher suites by setting PreferServerCipherSuites to true
	//another setting we are manulpulating here is the CurvePreference, this allow
	//to specify the elliptic curves preferred during TLS handshake.
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:		  []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:     cfg.Addr,
		//limit the max header length to 0.5MB
		MaxHeaderBytes: 524288
		ErrorLog: app1.ErrorLog,
		Handler:  app.routes(),
		TLSConfig: tlsConfig,
		//adding idle, read, write timeouts to the server.
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app1.InfoLog.Printf("Starting server on %s", cfg.Addr)
	app1.InfoLog.Println("static directory location", cfg.StaticDir)
	app1.InfoLog.Println("Is Production", cfg.IsProd)
	
	// Use the ListenAndServeTLS() method to start the HTTPS server. We
	// pass in the paths to the TLS certificate and corresponding private key as
	// the two parameters.
	err1 := srv.ListenAndServeTLS("F:/snippetbox_git/goproject/tls/cert.pem", "F:/snippetbox_git/goproject/tls/key.pem")
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
