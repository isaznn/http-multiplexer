package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	mockMuxDelayPerIterationMs = 10
	mockRequestLimit = 3
	urlPerRequestLimit = 1
)

/*
* Mock Get return map[<url>]: ok
* with delay
*/
type mockService struct {}
func (s *mockService) Mux(ctx context.Context, urls []string) (map[string]string, error) {
	m := make(map[string]string, len(urls))

	for _, url := range urls {
		// imitation of work
		time.Sleep(mockMuxDelayPerIterationMs * time.Millisecond)

		m[url] = "ok"
	}

	return m, nil
}

func TestHandler_MuxerHandler(t *testing.T) {
	t.Run("wrong method", func(t *testing.T) {
		// arrange
		h := NewHandler(mockRequestLimit, urlPerRequestLimit, &mockService{})
		r := h.InitRouter()
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/muxer", nil)

		// act
		r.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status code %v but got %v", http.StatusMethodNotAllowed, rr.Code)
		}
	})

	t.Run("empty urls array", func(t *testing.T) {
		// arrange
		h := NewHandler(mockRequestLimit, urlPerRequestLimit, &mockService{})
		r := h.InitRouter()
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/muxer", bytes.NewReader([]byte(`{"urls":[]}`)))

		// act
		r.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusBadRequest || !strings.Contains(rr.Body.String(), emptyArrayErrorText) {
			t.Error("incorrect handling of an empty array")
		}
	})

	t.Run("exceeding the limit urls per request", func(t *testing.T) {
		// arrange
		h := NewHandler(mockRequestLimit, urlPerRequestLimit, &mockService{})
		r := h.InitRouter()
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/muxer", bytes.NewReader([]byte(`{"urls":["https://domain.com","https://another.domain.net"]}`)))

		// act
		r.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusBadRequest || !strings.Contains(rr.Body.String(), tooMuchUrlErrorText) {
			t.Error("incorrect handling of url limit exceeded")
		}
	})

	t.Run("invalid url", func(t *testing.T) {
		// arrange
		h := NewHandler(mockRequestLimit, urlPerRequestLimit, &mockService{})
		r := h.InitRouter()
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/muxer", bytes.NewReader([]byte(`{"urls": ["https:://domain.com"]}`)))

		// act
		r.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusBadRequest || !strings.Contains(rr.Body.String(), invalidUrlErrorText) {
			t.Error("incorrect handling of url limit exceeded")
		}
	})

	t.Run("exceeding the request limit", func(t *testing.T) {
		// arrange
		var (
			totalRequests = 5
			code200counter int32
			code500counter int32
			wg sync.WaitGroup
		)
		h := NewHandler(mockRequestLimit, urlPerRequestLimit, &mockService{})
		r := h.InitRouter()

		// act
		for i := 0; i < totalRequests; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				rr := httptest.NewRecorder()
				req:= httptest.NewRequest(http.MethodPost, "/muxer", bytes.NewReader([]byte(`{"urls":["https://domain.com"]}`)))
				r.ServeHTTP(rr, req)
				switch rr.Code {
				case 200:
					atomic.AddInt32(&code200counter, IncDelta)
				case 500:
					atomic.AddInt32(&code500counter, IncDelta)
				}
			}()
			time.Sleep(100000 * time.Nanosecond)
		}
		wg.Wait()

		// assert
		if code200counter != mockRequestLimit || code500counter != int32(totalRequests - mockRequestLimit) {
			t.Error("incorrect handling of the request limit")
		}
	})
}
