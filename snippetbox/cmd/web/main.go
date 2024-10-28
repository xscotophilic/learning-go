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

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"snippetbox.xscotophilic.art/internal/env"
	"snippetbox.xscotophilic.art/internal/models"
)

type application struct {
	enableDebugMode bool
	errorLog        *log.Logger
	infoLog         *log.Logger
	formDecoder     *form.Decoder
	templateCache   map[string]*template.Template
	sessionManager  *scs.SessionManager
	snippets        models.SnippetModelInterface
	users           models.UserModelInterface
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Port")
	dsn := flag.String("dsn", env.DBCreds, "MySQL data source name")
	enableDebugMode := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.LUTC)
	errorLog := log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		enableDebugMode: *enableDebugMode,
		errorLog:        errorLog,
		infoLog:         infoLog,
		formDecoder:     formDecoder,
		templateCache:   templateCache,
		sessionManager:  sessionManager,
		snippets:        &models.SnippetModel{DB: db},
		users:           &models.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     errorLog,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
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
