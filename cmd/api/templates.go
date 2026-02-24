package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/EDITH5607/PasteDash/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
	Form any
	Flash string
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
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err!=nil {
		return nil,err
	}

	//iterate through pages and get name and store in cache
	for _,page:= range pages {
		//strip the name from path eg: ./ui/html/pages/hello.html -> hello.html and make this name as key in cache..
		name := filepath.Base(page)

		// parsing base html file
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}
		
		// parsing all "partial file instead of making an slice"
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")

		// making the "page" parse(template map) with the following files
		ts, err = ts.ParseFiles(page)
		if err!=nil {
			return nil,err
		}

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

	}
}