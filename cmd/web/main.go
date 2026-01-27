package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/EDITH5607/PasteDash/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct{
	ErrLog *log.Logger
	InfoLog *log.Logger
	Snippet *models.SnippetModel
}


func main() {
	//getting http network address through Cli
	addr := flag.String("addr",":4000","HTTP Network Address")
	dsn := flag.String("dsn", "web:helloworld@/pastedash?parseTime=true", "MySQL data source name!!!")
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


	// initialzing an struct of custom logging
	app := &application{
		ErrLog: ErrLog,
		InfoLog: InfoLog,
		Snippet: &models.SnippetModel{DB:db},
	} 

	// initializing server with custom error logging.
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: ErrLog,
		Handler: app.routes(),
	}

	//server initialization and starting
	InfoLog.Printf("Starting Server on :127.0.0.1:%s",*addr)
	err = srv.ListenAndServe()
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
