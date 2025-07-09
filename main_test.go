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

func TestBasicAuthMiddleware_Success(t *testing.T) {
	username := "user"
	password := "pass"
	called := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	handler := basicAuthMiddleware(username, password, next)

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth(username, password)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !called {
		t.Error("Expected next handler to be called with correct credentials")
	}
	if rec.Result().StatusCode == http.StatusUnauthorized {
		t.Error("Did not expect unauthorized status with correct credentials")
	}
}

func TestBasicAuthMiddleware_Failure(t *testing.T) {
	username := "user"
	password := "pass"
	called := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	handler := basicAuthMiddleware(username, password, next)

	req := httptest.NewRequest("GET", "/", nil)
	// No credentials set
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if called {
		t.Error("Did not expect next handler to be called without credentials")
	}
	if rec.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected unauthorized status, got %d", rec.Result().StatusCode)
	}

	// Try with wrong credentials
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.SetBasicAuth("wrong", "creds")
	rec2 := httptest.NewRecorder()
	called = false

	handler.ServeHTTP(rec2, req2)

	if called {
		t.Error("Did not expect next handler to be called with wrong credentials")
	}
	if rec2.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected unauthorized status, got %d", rec2.Result().StatusCode)
	}
}
