package main

import "net/http"

func (app *application) routes() *http.ServeMux{
	// router
	mux := http.NewServeMux()
	
	//file server
	fs:= http.FileServer(http.Dir("./ui/static"))

	//routes 
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)
	mux.Handle("/static/",http.StripPrefix("/static",fs))
	return mux
}