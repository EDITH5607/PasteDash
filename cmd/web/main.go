package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	//getting http network address through Cli
	addr := flag.String("addr",":4000","HTTP Network Address")
	flag.Parse()
	// router
	mux := http.NewServeMux()

	//file server
	fs:= http.FileServer(http.Dir("./ui/static"))

	//routes 
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	mux.Handle("/static/",http.StripPrefix("/static",fs))

	//server initialization and starting
	log.Printf("Starting Server on :127.0.0.1:%s",*addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)

}
