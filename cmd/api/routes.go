package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler{
	// router
	mux := http.NewServeMux()
	
	//file server
	fs:= http.FileServer(http.Dir("./ui/static"))

	//routes 
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)
	mux.Handle("/static/",http.StripPrefix("/static",fs))

	// use alice of middleware chaining alice.new(m1,m2,m3)  request->m1->m2->m3
	standard := alice.New(app.logRequest,securityHeader)

	// .then for adding handlers
	return standard.Then(mux)

}