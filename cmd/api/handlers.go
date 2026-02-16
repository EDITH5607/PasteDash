package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"github.com/EDITH5607/PasteDash/internal/models"
	"github.com/julienschmidt/httprouter"
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

	title := "0 snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n - Kobayashi Issa"
	expires := 7
	id, err := app.Snippet.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet"))
} 
