package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

/*
	we make outer function to pass the next http.handler{next middleware} because the returning innerfunction cannot accept next as parameter like
	func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		 next.ServeHTTP(w,r)
	}
	then this will not be a http.Handler type
	it's actually nesting them. Because each middleware calls next.ServeHTTP(), the computer creates a "stack."
    -secureHeaders starts...
        -Inside secureHeaders, logRequest starts...
            -Inside logRequest, homeHandler starts...
            -homeHandler finishes.
        -logRequest finishes.
    -secureHeaders finishes.
This is why we say it's sequential, but also wrapped. The outer functions stay "active" while the inner functions are running.


The outerfn execute when the server is started and return the innerfn . this inner fn will execute when ever the request comes...
*/


func (app *application)logRequest(next http.Handler) http.Handler {
	fmt.Printf("logRequest middleware called\n\n")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		 next.ServeHTTP(w,r)
	})
}

func securityHeader(next http.Handler) http.Handler {
	fmt.Printf("security header  middleware called\n\n")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w,r)
		fmt.Print("called security header after middleware")
	}) 	
}

func (app *application) requiredAuthentication(next http.Handler) http.Handler {
	fmt.Println("authentication middleware called")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w,r,"/user/login",http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control","no-store")
		next.ServeHTTP(w,r)
	
	})
}




// the actual object that returned by the noSurf middleware is the csrfhandler and when we write any code before that is just configering that 
func noSurf(next http.Handler) http.Handler {
	fmt.Println("surf middleware called")
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: true,
	})
	return csrfHandler
}