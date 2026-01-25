package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)


func (app *application)home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.NotFound(w)
		return 
	}
	
	files := []string{
		"./ui/html/pages/home.html",
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",

	}
	ts,err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w,"base", nil)
	if err!=nil {
		app.ErrLog.Println(err.Error())
		app.serverError(w, err)
		} 
}

func (app *application)snippetView(w http.ResponseWriter, r *http.Request) {
	id,err := strconv.Atoi(r.URL.Query().Get("id"))
	if err!= nil || id < 1 {
		app.NotFound(w)
		return
	}
	fmt.Fprintf(w,"Display a specific snippet with ID : %d", id)

}

func (app *application)snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) // this way we can send tokens for verification.
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new Snippet !!"))
}
