package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"markperea.com/snippetbox/internal/models"
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

func (a *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	td := a.newTemplateData(r)
	td.Form = snippetCreateForm{Expires: 365}
	a.render(w, r, http.StatusOK, "create.html", td)
}

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
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
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("Content"),
		Expires:     expires,
		FieldErrors: make(map[string]string),
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7, or 365"
	}
	if len(form.FieldErrors) > 0 {
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

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
