package main

import (
	"errors"
	"net/http"
	"strconv"
	"github.com/EDITH5607/PasteDash/internal/models"
	"github.com/julienschmidt/httprouter"
	"fmt"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// displaying latest code ....
	snippets, err := app.Snippet.Latest()
	if err != nil {
		app.NotFound(w)
		return
	}

	// getting year and other data using newtemplatedata fun and adding snippets to the struct from db.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// rendering the cached template of home page
	app.render(w, http.StatusOK, data, "home.html")



}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))	
	if err != nil || id < 1 {
		app.NotFound(w)
		return
	}

	snippet, err := app.Snippet.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.NotFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// getting year and other data using newtemplatedata fun and adding snippet to the struct from db.
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w,http.StatusOK, data, "view.html")
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, data, "create.html")
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	//r.ParseForm() which adds any data in POST request bodies to the r.PostForm map
	err :=r.ParseForm()
	if err!=nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//r.PostForm.Get() method to retrieve the title and content from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	expires,err := strconv.Atoi(r.PostForm.Get("expires"))
	if err!=nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//calling db method to insert the snippet set.
	id, err :=app.Snippet.Insert(title,content,expires)
	if err!=nil {
		app.serverError(w,err)
		return
	}

	// redirect to show the newly added snippet.
	http.Redirect(w,r, fmt.Sprintf("/snippet/view/%d",id),http.StatusSeeOther)

} 
