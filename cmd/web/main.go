package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct{
	ErrLog *log.Logger
	InfoLog *log.Logger
}


func main() {
	//getting http network address through Cli

	addr := flag.String("addr",":4000","HTTP Network Address")
	flag.Parse()


	// temp files for logging error,Infolog 
	/* /tmp/error.log and /tmp/info.log are the system tmp file system parts if we use './tmp/info.log' we make file on the project directory*/
	/* if you want to see the system temp folder just use command "tail -f /tmp/error.log "*/

	// e,eerr := os.OpenFile("/tmp/error.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	// e,eerr := os.OpenFile("/tmp/error.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	// if eerr != nil  {
	// 	log.Fatal(eerr)
	// }
	// i,ier := os.OpenFile("/tmp/info.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	// if ier != nil  {
	// 	log.Fatal(ier)
	// }
	// defer i.Close()
	// defer e.Close()



	//InfoLog and ErrorLog

	InfoLog := log.New(os.Stdout,"INFO\t", log.Ldate | log.Ltime)
	ErrLog := log.New(os.Stderr, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile)
	// initialzing an struct of custom logging
	app := &application{
		ErrLog: ErrLog,
		InfoLog: InfoLog,
	} 

	// initializing server with custom error logging.

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: ErrLog,
		Handler: app.routes(),
	}


	//server initialization and starting

	InfoLog.Printf("Starting Server on :127.0.0.1:%s",*addr)
	err := srv.ListenAndServe()
	ErrLog.Fatal(err)

}
