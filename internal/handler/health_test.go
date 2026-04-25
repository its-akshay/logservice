package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/logservice/internal/handler"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		c *gin.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest("GET", "/healthz", nil)
			c.Request = req

			handler.HealthHandler(tt.c)

			if w.Code != http.StatusOK {
				t.Errorf("expected 200, got %d", w.Code)
			}
		})
	}
}
