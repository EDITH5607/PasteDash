package main

import "net/http"


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
*/


func (app *application)logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		 next.ServeHTTP(w,r)
	})
}

func securityHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w,r)
	}) 	
}

func (app *application) requiredAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w,r,"/user/login",http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control","no-store")
		next.ServeHTTP(w,r)
	
	})
}
