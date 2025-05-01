package main

import "net/http"

func (a *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.Handle("GET /{$}", a.sessionManager.LoadAndSave(http.HandlerFunc(a.home)))

	mux.Handle("GET /snippet/view/{id}", a.sessionManager.LoadAndSave(http.HandlerFunc(a.snippetView)))
	mux.Handle("GET /snippet/create", a.sessionManager.LoadAndSave(http.HandlerFunc(a.snippetCreate)))
	mux.Handle("POST /snippet/create", a.sessionManager.LoadAndSave(http.HandlerFunc(a.snippetCreatePost)))

	mux.Handle("GET /user/signup", a.sessionManager.LoadAndSave(http.HandlerFunc(a.userSignup)))
	mux.Handle("POST /user/signup", a.sessionManager.LoadAndSave(http.HandlerFunc(a.userSignupPost)))
	mux.Handle("GET /user/login", a.sessionManager.LoadAndSave(http.HandlerFunc(a.userLogin)))
	mux.Handle("POST /user/login", a.sessionManager.LoadAndSave(http.HandlerFunc(a.userLoginPost)))
	mux.Handle("POST /user/logout", a.sessionManager.LoadAndSave(http.HandlerFunc(a.userLogoutPost)))

	return a.recoverPanic(a.logRequest(commonHeaders(mux)))
}
