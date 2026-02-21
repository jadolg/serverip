package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlers(t *testing.T) {
	tmpl, err := template.ParseFS(tmplFS, "index.html.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	mux := setupRoutes(tmpl)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	oldURL := wtfismyipIPv4URL
	wtfismyipIPv4URL = mockServer.URL
	defer func() { wtfismyipIPv4URL = oldURL }()

	t.Run("Health Check", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), "OK") {
			t.Errorf("expected body to contain OK, got %s", w.Body.String())
		}
	})

	t.Run("Home Page HTML", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		if w.Header().Get("Content-Type") != "text/html" {
			t.Errorf("expected content type text/html, got %s", w.Header().Get("Content-Type"))
		}
		if !strings.Contains(w.Body.String(), mockResponse.YourFuckingIPAddress) {
			t.Errorf("expected body to contain IP %s", mockResponse.YourFuckingIPAddress)
		}
	})

	t.Run("Home Page JSON", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("expected content type application/json, got %s", w.Header().Get("Content-Type"))
		}

		var got WTFIsMyIPData
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}
		if got.YourFuckingIPAddress != mockResponse.YourFuckingIPAddress {
			t.Errorf("expected IP %s, got %s", mockResponse.YourFuckingIPAddress, got.YourFuckingIPAddress)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/not-found", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})

	t.Run("External API Failure", func(t *testing.T) {
		// Point to an invalid URL to trigger failure
		wtfismyipIPv4URL = "http://invalid-url-that-does-not-exist"

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}
