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

	return a.recoverPanic(a.logRequest(commonHeaders(mux)))
}
