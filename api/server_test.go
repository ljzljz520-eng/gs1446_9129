package api

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealth(t *testing.T) {
	r := httptest.NewRecorder()
	HealthHandler(time.Now())(r, httptest.NewRequest("GET", "/health", nil))
	if r.Code != 200 {
		t.Fatal(r.Code)
	}
}
