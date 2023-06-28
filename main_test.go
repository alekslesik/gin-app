package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDefaultRoute(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()

	ctx, router := gin.CreateTestContext(w)
	setupRouter(router)

	req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	router.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		t.Fatalf("expected response code %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()

	expected := "Hello, gin!"

	if expected != strings.Trim(body, " \r\n") {
		t.Fatalf("expected response body '%s', got '%s'", expected, body)
	}

} // End TestDefaultRoute

