package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_MuxerHandler(t *testing.T) {
	t.Run("wrong method", func(t *testing.T) {
		// arrange
		h := NewHandler()
		r := h.InitRouter()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/muxer", nil)

		// act
		r.ServeHTTP(rec, req)

		// assert
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status code %v but got %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}
