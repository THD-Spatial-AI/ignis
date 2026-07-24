package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thd-spatial-ai/ignis/internal/api/handler"

	"github.com/gin-gonic/gin"
)

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", HealthCheck)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != `{"status":"OK"}` {
		t.Errorf("body = %q, want %q", w.Body.String(), `{"status":"OK"}`)
	}
}

func TestRegisterRoutes_registersExpectedRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// handler.New stores the pool without dialing it, so a nil pool is safe
	// here - route registration never issues a query.
	h := handler.New(nil, "tabula")

	RegisterRoutes(r, h)

	want := map[string]bool{
		"GET /favicon.ico":                         false,
		"GET /health":                              false,
		"GET /api/v1/data/:code":                   false,
		"GET /api/v1/variants/:country_iso2":       false,
		"GET /api/v1/variants/:country_iso2/match": false,
		"GET /api/v1/fields":                       false,
		"POST /api/v1/calculate/:code":             false,
	}

	for _, route := range r.Routes() {
		key := route.Method + " " + route.Path
		if _, expected := want[key]; expected {
			want[key] = true
		}
	}

	for route, found := range want {
		if !found {
			t.Errorf("expected route %q to be registered", route)
		}
	}
}

func TestRegisterRoutes_faviconReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.New(nil, "tabula")
	RegisterRoutes(r, h)

	req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}
