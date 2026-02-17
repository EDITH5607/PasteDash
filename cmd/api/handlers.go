package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
	
	"github.com/EDITH5607/PasteDash/internal/models"
	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title string
	Content string
	Expires int
	fieldErrors map[string]string

}



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

	// extracting params from the request context
	params := httprouter.ParamsFromContext(r.Context())

	// convert string to int
	id, err := strconv.Atoi(params.ByName("id"))	
	if err != nil || id < 1 {
		app.NotFound(w)
		return
	}

	// getting the snippet using the id from db.
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

	expires,err := strconv.Atoi(r.PostForm.Get("expires"))
	if err!=nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//r.PostForm.Get() method to retrieve the title and content from the r.PostForm map.
	// use snippetCreateform for hold form data and empty map for any validation errors.if data is incorrect the other form data is stored and user only need to change the wrong one.
	form := &snippetCreateForm{
		Title : r.PostForm.Get("title"),
		Content : r.PostForm.Get("content"),
		Expires: expires,
		fieldErrors: map[string]string{},


	}

	// validation of the data
	if strings.TrimSpace(form.Title) == "" {
		form.fieldErrors["title"] = "This field cannot be blank"

	// utf8.RuneCountInString used intead of len fun because of different characters are not accounted eg:ë
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.fieldErrors["content"] = "This field cannot be blank"
	}

	if  form.Expires !=1 && form.Expires !=7 && form.Expires != 365 {
		form.fieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	// if any error in the map just send reponse the error and return the handler
	if len(form.fieldErrors) > 0 {
		// passing in the previous data and validation errors
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, data, "create.html")
		return
	}

	//calling db method to insert the snippet set.
	id, err :=app.Snippet.Insert(form.Title,form.Content,form.Expires)
	if err!=nil {
		app.serverError(w,err)
		return
	}

	// redirect to show the newly added snippet.
	http.Redirect(w,r, fmt.Sprintf("/snippet/view/%d",id),http.StatusSeeOther)

} 
