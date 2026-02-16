package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler{
	// router
	router := httprouter.New()	

	// custom handler to show notfound error instead of the router's own error handling
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				app.NotFound(w)
	}) 

	//file server
	fs:= http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet,"/static/*filepath",http.StripPrefix("/static",fs))


	//routes 
	router.HandlerFunc(http.MethodGet,"/", app.home)
	router.HandlerFunc(http.MethodGet,"/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet,"/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost,"/snippet/create", app.snippetCreatePost)

	// use alice of middleware chaining alice.new(m1,m2,m3)  request->m1->m2->m3
	standard := alice.New(app.logRequest,securityHeader)

	// .then for adding handlers
	return standard.Then(router)

}