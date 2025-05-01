package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"markperea.com/snippetbox/internal/models"
	"markperea.com/snippetbox/internal/validator"
)

func (a *application) home(w http.ResponseWriter, r *http.Request) {
	ss, err := a.snippets.Latest()
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	td := a.newTemplateData(r)
	td.Snippets = ss
	a.render(w, r, http.StatusOK, "home.html", td)
}

func (a *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	s, err := a.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			a.serverError(w, r, err)
		}
	}

	td := a.newTemplateData(r)
	td.Snippet = s
	a.render(w, r, http.StatusOK, "view.html", td)
}

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

func (a *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	td := a.newTemplateData(r)
	td.Form = snippetCreateForm{Expires: 365}
	a.render(w, r, http.StatusOK, "create.html", td)
}

func (a *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")

	if !form.Valid() {
		td := a.newTemplateData(r)
		td.Form = form
		a.render(w, r, http.StatusUnprocessableEntity, "create.html", td)
		return
	}

	id, err := a.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	a.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

type userSignupForm struct {
	Name     string
	Email    string
	Password string
	validator.Validator
}

func (a *application) userSignup(w http.ResponseWriter, r *http.Request) {
	td := a.newTemplateData(r)
	td.Form = userSignupForm{}
	a.render(w, r, http.StatusOK, "signup.html", td)
}

func (a *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form := userSignupForm{
		Name:     r.PostForm.Get("name"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		td := a.newTemplateData(r)
		td.Form = form
		a.render(w, r, http.StatusUnprocessableEntity, "signup.html", td)
	}

	err = a.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			td := a.newTemplateData(r)
			td.Form = form
			a.render(w, r, http.StatusUnprocessableEntity, "signup.html", td)
		} else {
			a.serverError(w, r, err)
		}
		return
	}

	a.sessionManager.Put(r.Context(), "flash", "Your sign up was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email    string
	Password string
	validator.Validator
}

func (a *application) userLogin(w http.ResponseWriter, r *http.Request) {
	td := a.newTemplateData(r)
	td.Form = userLoginForm{}
	a.render(w, r, http.StatusOK, "login.html", td)

}

func (a *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form := userLoginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		td := a.newTemplateData(r)
		td.Form = form
		a.render(w, r, http.StatusUnprocessableEntity, "login.html", td)
	}

	id, err := a.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			td := a.newTemplateData(r)
			td.Form = form
			a.render(w, r, http.StatusUnprocessableEntity, "login.html", td)
		} else {
			a.serverError(w, r, err)
		}
		return
	}

	if err := a.sessionManager.RenewToken(r.Context()); err != nil {
		a.serverError(w, r, err)
		return
	}

	a.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (a *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	if err := a.sessionManager.RenewToken(r.Context()); err != nil {
		a.serverError(w, r, err)
		return
	}

	a.sessionManager.Remove(r.Context(), "authenticatedUserID")

	a.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
