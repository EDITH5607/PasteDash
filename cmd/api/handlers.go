package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/EDITH5607/PasteDash/internal/models"
	"github.com/EDITH5607/PasteDash/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title string `form:"title"`
	Content string `form:"content"`
	Expires int  `form:"expires"`
	validator.Validator  `form:"-"`
}

type userSignupForm struct {
	Name 		string   `form:"name"`
	Email		string   `form:"email"`
	Password 	string   `form:"password"`
	validator.Validator  `form:"-"`

}

type userLoginForm struct {
	Email string `form:"email"`
	Password string `form:"password"`
	validator.Validator `form:"-"`
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
	data.Form = &snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, data, "create.html")
}




func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	
	var form snippetCreateForm

	// use decoderpostform to parse the data from the form and show if any invalid encode error happens
	err := app.DecodePostForm(r, &form)
	if err != nil {
		app.clientError(w,http.StatusBadRequest)
		return
	}

	// validation of the data
	form.CheckField(validator.NotBlank(form.Title), "title",  "This field cannot be blank")
	form.CheckField(validator.MaxChar(form.Title,100), "title","This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermitInt(form.Expires, 1,7,365),"expires", "This field must equal 1, 7 or 365")

	// if any error in the map form.valid return false!!
	if  !form.Valid() {
		// making a template data struct
		data := app.newTemplateData(r)
		// passing snippetform to templatedata,to provide previous data and to solve validation errors
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

	// sending the message "flash":"Snippet Sucessfully Created!!" to the current request context
	app.sessionManager.Put(r.Context(), "flash","Snippet Successfully Created!!")

	// redirect to show the newly added snippet.
	http.Redirect(w,r, fmt.Sprintf("/snippet/view/%d",id),http.StatusSeeOther)

} 



//signup GET function for rendering html page.
func (app *application) userSignup (w http.ResponseWriter, r *http.Request){
	data := app.newTemplateData(r)
	data.Form = &userSignupForm{}
	app.render(w, http.StatusOK, data, "signup.html")
}


// signup POST function for validation and injection data to the DB
func (app *application) userSignupPost (w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	// this will decode the form and parse it in to the form(usesignupform)
	err := app.DecodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be Blank !!")
	form.CheckField(validator.NotBlank(form.Email),"email", "This field cannot be Blank !!")
	form.CheckField(validator.Matches(form.Email),"email", "Please Enter valid Email !!")
	form.CheckField(validator.NotBlank(form.Password), "password", "This Field cannot be blank !!")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This Field must be atleast 8 character long !!")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w,http.StatusUnprocessableEntity, data,"signup.html")
		return
	}

	err = app.Users.Insert(form.Name, form.Email, form.Password)
	if err!= nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity,data, "signup.html")
		} else {
			app.serverError(w,err)
		}
		return
	}
	app.sessionManager.Put(r.Context(), "flash","Your signup was successful. Please log in." )
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app * application) userLogin (w http.ResponseWriter, r *http.Request)  {
	data := app.newTemplateData(r)
	data.Form = &userLoginForm{}
	app.render(w,http.StatusOK, data, "login.html")
}

func (app *application) userLoginPost (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user")
}

func (app *application) userLogoutPost (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}