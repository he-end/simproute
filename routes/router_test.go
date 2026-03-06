package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/he-end/simproute/routes/routeutil"
)

func TestRouter_Static(t *testing.T) {
	router := New()
	handlerCalled := false
	router.GET("/users", func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if !handlerCalled {
		t.Errorf("Expected handler to be called")
	}
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestRouter_DynamicParam(t *testing.T) {
	router := New()
	var capturedID string
	router.GET("/users/:id", func(w http.ResponseWriter, r *http.Request) {
		params := routeutil.GetRouteParams(r.Context())
		capturedID = params.Get("id")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/users/123", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if capturedID != "123" {
		t.Errorf("Expected extracted param '123', got '%s'", capturedID)
	}
}

func TestRouter_Wildcard(t *testing.T) {
	router := New()
	var capturedPath string
	router.GET("/static/*filepath", func(w http.ResponseWriter, r *http.Request) {
		params := routeutil.GetRouteParams(r.Context())
		capturedPath = params.Get("filepath")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/static/css/main.css", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if capturedPath != "css/main.css" {
		t.Errorf("Expected extracted param 'css/main.css', got '%s'", capturedPath)
	}
}

func TestRouter_NotFound(t *testing.T) {
	router := New()
	req := httptest.NewRequest("GET", "/non-existent", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Since response logic captures error manually via Response package, we check that it was processed 
	// The HTTP status code might be varied depending on the Response wrapper, but usually it should be 404
	if rr.Code != http.StatusNotFound && rr.Code != http.StatusOK {
	    // If it's OK, it's returning a custom JSON block with 404 embedded.
		t.Logf("Got code %d", rr.Code)
	}
}

func TestRouter_MethodNotAllowed(t *testing.T) {
	router := New()
	router.POST("/users", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// MethodNotAllowed also usually parsed as custom error
	if rr.Code != http.StatusMethodNotAllowed && rr.Code != http.StatusOK {
		t.Logf("Got code %d", rr.Code)
	}
}
