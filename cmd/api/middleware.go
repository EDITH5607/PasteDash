package main

import (
	"context"
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

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// logging each informtion about the request...
		app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func securityHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// setting the security header...
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (app *application) requiredAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check the useauthenticationkey in the request and give permission to access the private routes  
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// tell the browser not to store the html pages...
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)

	})
}

// the actual object that returned by the noSurf middleware is the csrfhandler and when we write any code before that is just configering that
func noSurf(next http.Handler) http.Handler {
	// making a new handler and set the csrf token cookie
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	// returning the configered csrf handler
	return csrfHandler
}

//authenticate middleware for checking the userauthenication id from cookie(request) and db id of the user

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// getting the authenticationuserid (as integer) from the request context
		id := app.sessionManager.GetInt(r.Context(), "AuthenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// checking is the same id exists in db
		exists, err := app.Users.Exists(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// if exists then make a new child context with key:value as shown
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)

	})
}
