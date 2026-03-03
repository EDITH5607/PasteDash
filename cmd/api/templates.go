package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"github.com/EDITH5607/PasteDash/internal/models"
	"github.com/EDITH5607/PasteDash/ui"
	"github.com/justinas/nosurf"
)

type templateData struct {
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
	Form any
	Flash string
	IsAuthenticated  bool
	CSRFToken string

}

func humanDate(t time.Time) string {
	return  t.Format("02 Jan 2006 at 15:06")
}

var functions  = template.FuncMap{
	"humanDate":humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	//making a map as key(string): value(template.Template)
	cache := map[string]*template.Template{}

	// gets the .html files from the patturn.
	// pages, err := filepath.Glob("./ui/html/pages/*.html")
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err!=nil {
		return nil,err
	}

	//iterate through pages and get name and store in cache
	for _,page:= range pages {
		//strip the name from path eg: ./ui/html/pages/hello.html -> hello.html and make this name as key in cache..
		name := filepath.Base(page)

		patterns := []string {
			"html/base.html", 
			"html/partials/*.html",
			page,
		}

		// parsing base html file
		// ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		ts , err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
				
		// // parsing all "partial file instead of making an slice"
		// ts, err = ts.ParseGlob("./ui/html/partials/*.html")

		// // making the "page" parse(template map) with the following files
		// ts, err = ts.ParseFiles(page)
		// if err!=nil {
		// 	return nil,err
		// }

		// making the page name as key and template object of correponding as key. like hello.html: template of hello.html
		cache[name] = ts
	}
	return cache, nil
}

func (app *application) newTemplateData(r *http.Request) *templateData{
	return &templateData{
		CurrentYear: time.Now().Year(),
		// popstring() will retrieve and remove it from the session data...
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
		// check the isauthentication helper for rendering the create in navbar and other html pages
		IsAuthenticated: app.isAuthenticated(r),
		// generating the csrf token to the template placeholders
		CSRFToken: nosurf.Token(r),

	}
}
