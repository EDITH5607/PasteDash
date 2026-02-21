package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/EDITH5607/PasteDash/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct{
	ErrLog *log.Logger
	InfoLog *log.Logger
	Snippet *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder *form.Decoder
	sessionManager *scs.SessionManager
}


func main() {
	//getting http network address through Cli
	addr := flag.String("addr",":4000","HTTP Network Address")
	dsn := flag.String("dsn", "web:#354286Aatt@/pastedash?parseTime=true", "MySQL data source name!!!")
	flag.Parse()


	//InfoLog and ErrorLog
	InfoLog := log.New(os.Stdout,"INFO\t", log.Ldate | log.Ltime)
	ErrLog := log.New(os.Stderr, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile)


	//db connection
	db,err := openDB(*dsn)
	if err!=nil {
		ErrLog.Fatal(err)	
	}
	defer db.Close()

	//template cache..
	templatecache,err := newTemplateCache()
	if err!=nil {
		ErrLog.Fatalln(err)
	}

	// initialize form decoder
	formDecoder := form.NewDecoder()


	// creating an instance of session Manager and passing the mysql db to store the session and set the lifetime.
	sessionManager := scs.New()
	// 'sessions' name is default name in session manager package for storing user sessions. so we made a table called sessions in our db
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 *time.Hour
	sessionManager.Cookie.Secure = true



	// initialzing an struct of custom logging and passing db connection...(dependency injection)
	app := &application{
		ErrLog: ErrLog,
		InfoLog: InfoLog,
		Snippet: &models.SnippetModel{DB:db},
		templateCache: templatecache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	} 
	

	// initializing server with custom error logging.
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: ErrLog,
		Handler: app.routes(),
	}

	//server initialization and starting
	InfoLog.Printf("Starting Server on :127.0.0.1:%s",*addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	ErrLog.Fatal(err)

}

func openDB(dsn string) (*sql.DB, error) {
	db,err := sql.Open("mysql",dsn)
	if err!=nil {
		return nil,err
	}
	if err = db.Ping(); err!=nil {
		return nil,err
	}
	fmt.Println("DB Successfully Connected....")
	return db,nil
}
