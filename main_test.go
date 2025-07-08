package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	called := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	handler := loggingMiddleware(next)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !called {
		t.Error("Expected next handler to be called")
	}
}

func TestFileServer(t *testing.T) {
	dir := t.TempDir()

	content := "hello"
	filename := dir + "/hello.txt"
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fs := http.FileServer(http.Dir(dir))
	handler := loggingMiddleware(fs)

	req := httptest.NewRequest("GET", "/hello.txt", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != content {
		t.Errorf("expected body %q, got %q", content, string(body))
	}
}
