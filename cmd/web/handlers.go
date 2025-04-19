package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"markperea.com/snippetbox/internal/models"
)

func (a *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

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
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func (a *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// TODO: Remove test data
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := a.snippets.Insert(title, content, expires)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
