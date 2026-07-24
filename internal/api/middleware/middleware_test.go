package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestEngine(handlers ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(handlers...)
	r.POST("/echo", func(c *gin.Context) {
		body := make([]byte, 0)
		buf := make([]byte, 512)
		for {
			n, err := c.Request.Body.Read(buf)
			body = append(body, buf[:n]...)
			if err != nil {
				break
			}
		}
		c.String(http.StatusOK, "%d", len(body))
	})
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	return r
}

func TestRequestBodyLimit_allowsSmallBody(t *testing.T) {
	r := newTestEngine(RequestBodyLimit())
	body := bytes.Repeat([]byte("a"), 100)
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewReader(body))
	req.ContentLength = int64(len(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRequestBodyLimit_rejectsOversizedBody(t *testing.T) {
	r := newTestEngine(RequestBodyLimit())
	body := bytes.Repeat([]byte("a"), maxRequestBodyBytes+1)
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewReader(body))
	req.ContentLength = int64(len(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want %d", w.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestRequestLogger_passesThroughAndLogsStatus(t *testing.T) {
	r := newTestEngine(RequestLogger())
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != "pong" {
		t.Errorf("body = %q, want %q", w.Body.String(), "pong")
	}
}

func TestCORS_allowedOrigin(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:5173,https://app.example.com")
	r := newTestEngine(CORS())

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:5173")
	}
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCORS_disallowedOrigin(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "https://app.example.com")
	r := newTestEngine(CORS())

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Access-Control-Allow-Origin = %q, want empty for disallowed origin", got)
	}
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d (request still proceeds, just without CORS headers)", w.Code, http.StatusOK)
	}
}

func TestCORS_noOriginUnsetEnv(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "")
	r := newTestEngine(CORS())

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Access-Control-Allow-Origin = %q, want empty", got)
	}
}

func TestCORS_optionsPreflight_isAborted(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "https://app.example.com")
	r := newTestEngine(CORS())

	req := httptest.NewRequest(http.MethodOptions, "/ping", nil)
	req.Header.Set("Origin", "https://app.example.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
	if strings.TrimSpace(w.Body.String()) != "" {
		t.Errorf("body = %q, want empty (handler must not run after preflight abort)", w.Body.String())
	}
}
