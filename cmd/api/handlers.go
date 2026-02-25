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


// struct for snippet form
type snippetCreateForm struct {
	Title string `form:"title"`
	Content string `form:"content"`
	Expires int  `form:"expires"`
	validator.Validator  `form:"-"`
}

//struct for signup form
type userSignupForm struct {
	Name 		string   `form:"name"`
	Email		string   `form:"email"`
	Password 	string   `form:"password"`
	validator.Validator  `form:"-"`

}

// struct for login form
type userLoginForm struct {
	Email string `form:"email"`
	Password string `form:"password"`
	validator.Validator `form:"-"`
}



// Home page handler
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




// snippet view handler GET
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




// Snippet create page GET Handler
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	// to start with default value as expire:365
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
		// getting form field error messages  and other data using newtemplatedata func and adding snippet to the struct from db.
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
	// getting form field error messages  and other data using newtemplatedata func and adding snippet to the struct from db.
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
	// validationing the form data 
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be Blank !!")
	form.CheckField(validator.NotBlank(form.Email),"email", "This field cannot be Blank !!")
	form.CheckField(validator.Matches(form.Email),"email", "Please Enter valid Email !!")
	form.CheckField(validator.NotBlank(form.Password), "password", "This Field cannot be blank !!")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This Field must be atleast 8 character long !!")

	// flashing error if anything in the field error slice
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		// getting form field error messages  and other data using newtemplatedata func and adding snippet to the struct from db.
		app.render(w,http.StatusUnprocessableEntity, data,"signup.html")
		return
	}

	// insert the form data to the db (user)
	err = app.Users.Insert(form.Name, form.Email, form.Password)
	if err!= nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			// getting form field error messages  and other data using newtemplatedata func and adding snippet to the struct from db.
			app.render(w, http.StatusUnprocessableEntity,data, "signup.html")
		} else {
			app.serverError(w,err)
		}
		return
	}
	// flashing the successful message
	app.sessionManager.Put(r.Context(), "flash","Your signup was successful. Please log in." )
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// login GET handler for rendering pages
func (app * application) userLogin (w http.ResponseWriter, r *http.Request)  {
	data := app.newTemplateData(r)
	data.Form = &userLoginForm{}
	// getting form field error messages  and other data using newtemplatedata func and adding snippet to the struct from db.
	app.render(w,http.StatusOK, data, "login.html")
}

// login POST handler for form 
func (app *application) userLoginPost (w http.ResponseWriter, r *http.Request) {
	//making a zero instance and store data from form to the struct
	var form userLoginForm
	err := app.DecodePostForm(r,&form)
	if err!=nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validating form data 
	form.CheckField(validator.NotBlank(form.Email), "email", "This Field Cannot be Blank !!")
	form.CheckField(validator.Matches(form.Email), "email", "Please Enter valid Email !!")
	form.CheckField(validator.NotBlank(form.Password), "password", "This Field cannot be blank !!")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		// getting form field error messages  and other data using newtemplatedata func and adding snippet to the struct from db.
		app.render(w,http.StatusOK, data,"login.html")
	}

	// Authenticating the user with email and password
	id, err :=app.Users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err,models.ErrInvalidCredential) {
			form.AddNonFieldErrors("Email or password is incorrect!!")
			data := app.newTemplateData(r)
			data.Form = form
			// getting form field error messages  and other data using newtemplatedata func and adding snippet to the struct from db.
			app.render(w,http.StatusUnprocessableEntity, data, "login.html")
		} else {
			app.serverError(w,err)
		}
		return
	}

	// renew the token for security purpose to avoid unnecessary cyber attacks , the middleware will give a new token if token is not present when we enter the route so for changing that 
	err = app.sessionManager.RenewToken(r.Context())
	if err!=nil {
		app.serverError(w, err)
		return
	}

	// storing the authenticated user id as context for future use
	app.sessionManager.Put(r.Context(), "AuthenticatedUserID", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}	

// logout POST handler for logout and session cleaning
func (app *application) userLogoutPost (w http.ResponseWriter, r *http.Request) {
	// for security renew the token , the middleware will give a new token if token is not present when we enter the route so for changing that 
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w,err)
		return
	}
	// remove the authentication id from the context
	app.sessionManager.Remove(r.Context(), "AuthenticatedUserID")
	// showing the flash message
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w,r,"/", http.StatusSeeOther)

}