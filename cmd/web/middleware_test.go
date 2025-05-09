package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"markperea.com/snippetbox/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	commonHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	headers := []struct {
		header string
		value  string
	}{
		{
			header: "Content-Security-Policy",
			value:  "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		},
		{
			header: "Referrer-Policy",
			value:  "origin-when-cross-origin",
		},
		{
			header: "X-Content-Type-Options",
			value:  "nosniff",
		},
		{
			header: "X-Frame-Options",
			value:  "deny",
		},
		{
			header: "X-XSS-Protection",
			value:  "0",
		},
		{
			header: "Server",
			value:  "Go",
		},
	}

	for _, hh := range headers {
		assert.Equal(t, rs.Header.Get(hh.header), hh.value)
	}

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
