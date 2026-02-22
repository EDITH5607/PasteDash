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


	dynamic := alice.New(app.sessionManager.LoadAndSave)

	//routes 
	// middleware like loadAndSave accept the next fn as handler so we convert it http.handler type
	//also middleware return http.Handler type but the router.HandlerFunc accepts http.HandlerFunc but middleware return http.handler method that why we either convert to handler or change the router.Hander() 
	router.Handler(http.MethodGet,"/", app.sessionManager.LoadAndSave(http.HandlerFunc(app.home)))
	
	// the router.HandlerFunc() is remove because it accept fn and type convert the fn into http handler 
	// but the middleware chain and middleware are already type as handlers so no need to convert so we use router.Handler 
	router.Handler(http.MethodGet,"/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet,"/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost,"/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))


	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout",dynamic.ThenFunc(app.userLogoutPost))

	// use alice of middleware chaining alice.new(m1,m2,m3)  request->m1->m2->m3
	standard := alice.New(app.logRequest,securityHeader)

	// .then for adding handlers
	return standard.Then(router)

}