package main

import (
	"net/http"

	"markperea.com/snippetbox/ui"

	"github.com/justinas/alice"
)

func (a *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /ping", ping)

	dynamic := alice.New(a.sessionManager.LoadAndSave, noSurf, a.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(a.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(a.snippetView))

	protected := dynamic.Append(a.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(a.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(a.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(a.userLogoutPost))

	loggedOutOnly := dynamic.Append(a.requireUnauthenticated)

	mux.Handle("GET /user/signup", loggedOutOnly.ThenFunc(a.userSignup))
	mux.Handle("POST /user/signup", loggedOutOnly.ThenFunc(a.userSignupPost))
	mux.Handle("GET /user/login", loggedOutOnly.ThenFunc(a.userLogin))
	mux.Handle("POST /user/login", loggedOutOnly.ThenFunc(a.userLoginPost))

	standard := alice.New(a.recoverPanic, a.logRequest, commonHeaders)

	return standard.Then(mux)
}
