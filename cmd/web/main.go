package main

import (
	"log"
	"net/http"
)

func main() {
	// router
	mux := http.NewServeMux()

	//routes 
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	//server initialization and starting
	log.Println("Starting Server on :127.0.0.1:4000")
	err := http.ListenAndServe("127.0.0.1:4000", mux)
	log.Fatal(err)

}
