package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"github.com/go-playground/form/v4"
)


func (app *application)render(w http.ResponseWriter, status int, data *templateData, page string) {
	
	// passing name of the page for retrive from the cache if not found error handling works
	ts,ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w,err)
		return 
	}

	//initialzing new byte buffer
	buf := new(bytes.Buffer)

	//executing the template with its data, inital template name,  and Write the template to the buffer, instead of straight to the http.ResponseWriter. If there's an error the html page is not send to browser
	err := ts.ExecuteTemplate(buf,"base", data)
	if err !=nil {
		app.serverError(w,err)	
		return
	}
	// adding the status header like 200,405 as the header.
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWrite
	// another method instead of w.write() to reponse....
	buf.WriteTo(w)
	

}


func  (app *application)DecodePostForm(r *http.Request, dst any) error  {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	// use decoder to parse the form pass to dst(d)
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidEncodeError
		// errors.As() used to check the error chain if match found then copy to invaliddecder error variable that why we pass pointer.
		if errors.As(err,&invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return  nil
	
}


func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) NotFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) isAuthenticated (r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	 return isAuthenticated
}
