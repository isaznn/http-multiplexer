package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

const (
	mockRequestLimit            = 3
	mockConcurrentRequestsLimit = 3
)

type mockService struct {
	concurrentRequestsLimit int
}

type mockSafeMap struct {
	M map[string]string
	sync.Mutex
}

func newSafeMap(l int) *mockSafeMap {
	return &mockSafeMap{
		M: make(map[string]string, l),
	}
}

func (s *mockSafeMap) AllEntities() map[string]string {
	s.Lock()
	defer s.Unlock()
	return s.M
}

func (s *mockSafeMap) Store(key string, value string) {
	s.Lock()
	defer s.Unlock()
	s.M[key] = value
}

func (s *mockService) chunks(urls []string) [][]string {
	var dividedUrls [][]string
	chunkSize := (len(urls) + s.concurrentRequestsLimit - 1) / s.concurrentRequestsLimit

	for i := 0; i < len(urls); i += chunkSize {
		end := i + chunkSize
		if end > len(urls) {
			end = len(urls)
		}

		dividedUrls = append(dividedUrls, urls[i:end])
	}

	return dividedUrls
}

func (s *mockService) Mux(ctx context.Context, urls []string) (map[string]string, error) {
	m := newSafeMap(len(urls))
	chunks := s.chunks(urls)
	wrpCtx, cancel := context.WithCancel(ctx)
	errCh := make(chan struct{})
	wg := sync.WaitGroup{}
	isError := false

	// if received error - cancel context
	go func() {
		<-errCh
		isError = true
		cancel()
	}()

	// push chunks to goroutines
	for _, chunk := range chunks {
		wg.Add(1)
		go func(urls []string) {
			defer wg.Done()
			for _, url := range urls {
				select {
				case <-wrpCtx.Done():
					return
				default:
					httpReq, err := http.NewRequest(http.MethodGet, url, nil)
					if err != nil {
						errCh <- struct{}{}
						return
					}

					httpClient := http.Client{Timeout: 1 * time.Second}
					httpResp, err := httpClient.Do(httpReq)
					if err != nil {
						errCh <- struct{}{}
						return
					}
					defer httpResp.Body.Close()

					if httpResp.Status != "200 OK" {
						errCh <- struct{}{}
						return
					}

					respBytes, err := io.ReadAll(httpResp.Body)
					if err != nil {
						errCh <- struct{}{}
						return
					}
					m.Store(url, string(respBytes))
				}
			}
		}(chunk)
	}
	wg.Wait()

	switch isError {
	case true:
		return nil, fmt.Errorf("request ended with an error")
	default:
		return m.AllEntities(), nil
	}
}

func TestHandler_MuxerHandler(t *testing.T) {
	t.Run("wrong method", func(t *testing.T) {
		// arrange
		s := &mockService{
			concurrentRequestsLimit: mockConcurrentRequestsLimit,
		}
		h := NewHandler(mockRequestLimit, s)
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

	t.Run("correct processing urls", func(t *testing.T) {
		// arrange
		type testType struct {
			Value string
		}
		testData := testType{}
		testContent := "Hello World"
		testServerResponse := muxHandlerResponse{}
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(res).Encode(testType{Value: testContent})
			if err != nil {
				t.Errorf(err.Error())
			}
		}))
		defer func() {
			testServer.Close()
		}()
		s := &mockService{
			concurrentRequestsLimit: mockConcurrentRequestsLimit,
		}
		h := NewHandler(mockRequestLimit, s)
		r := h.InitRouter()
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/muxer",
			bytes.NewReader([]byte(fmt.Sprintf(`{"urls": ["%s"]}`, testServer.URL))),
		)

		// act
		r.ServeHTTP(rr, req)

		// assert
		err := json.NewDecoder(rr.Body).Decode(&testServerResponse)
		if err != nil {
			t.Errorf(err.Error())
		}
		if len(testServerResponse.Result) < 1 {
			t.Errorf("response length must be 1 but got %d", len(testServerResponse.Result))
		}
		err = json.Unmarshal([]byte(testServerResponse.Result[testServer.URL]), &testData)
		if err != nil {
			t.Errorf(err.Error())
		}
		if testData.Value != testContent {
			t.Errorf("expected value '%s' but got '%s'", testContent, testData.Value)
		}
	})

	t.Run("error in one of the urls", func(t *testing.T) {
		testServerOk := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte("ok"))
		}))
		defer func() {
			testServerOk.Close()
		}()
		testServerBad := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte("error"))
		}))
		defer func() {
			testServerBad.Close()
		}()
		s := &mockService{
			concurrentRequestsLimit: mockConcurrentRequestsLimit,
		}
		h := NewHandler(mockRequestLimit, s)
		r := h.InitRouter()
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/muxer",
			bytes.NewReader([]byte(fmt.Sprintf(`{"urls": ["%s", "%s"]}`, testServerOk.URL, testServerBad.URL))),
		)

		// act
		r.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code '%d' but got '%d'", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("exceeding the request limit", func(t *testing.T) {
		// arrange
		totalRequests := 5
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			time.Sleep(100 * time.Millisecond)
			res.WriteHeader(http.StatusOK)
			res.Write([]byte("ok"))
		}))
		defer func() {
			testServer.Close()
		}()
		s := &mockService{
			concurrentRequestsLimit: mockConcurrentRequestsLimit,
		}
		h := NewHandler(mockRequestLimit, s)
		r := h.InitRouter()
		var code200counter int
		var code500counter int
		var wg sync.WaitGroup

		// act
		for i := 0; i < totalRequests; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				rr := httptest.NewRecorder()
				req:= httptest.NewRequest(
					http.MethodPost,
					"/muxer",
					bytes.NewReader([]byte(fmt.Sprintf(`{"urls": ["%s"]}`, testServer.URL))),
				)
				r.ServeHTTP(rr, req)
				t.Log(rr.Code)
				t.Log(rr.Body.String())
				switch rr.Code {
				case 200:
					code200counter++
				case 500:
					code500counter++
				}
			}()
		}
		wg.Wait()

		// assert
		if code200counter != mockRequestLimit || code500counter != (totalRequests - mockRequestLimit) {
			t.Error("incorrect processing of the request limit")
		}
	})

	t.Run("concurrent processing urls", func(t *testing.T) {
		// arrange
		totalRequests := 3
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			time.Sleep(100 * time.Millisecond)
			res.WriteHeader(http.StatusOK)
			res.Write([]byte("ok"))
		}))
		defer func() {
			testServer.Close()
		}()
		s := &mockService{
			concurrentRequestsLimit: mockConcurrentRequestsLimit,
		}
		h := NewHandler(mockRequestLimit, s)
		r := h.InitRouter()
		now := time.Now()
		var wg sync.WaitGroup

		// act
		wg.Add(totalRequests)
		for i := 0; i < totalRequests; i++ {
			go func() {
				defer wg.Done()
				rr := httptest.NewRecorder()
				req:= httptest.NewRequest(
					http.MethodPost,
					"/muxer",
					bytes.NewReader([]byte(fmt.Sprintf(`{"urls": ["%s"]}`, testServer.URL))),
				)
				r.ServeHTTP(rr, req)
			}()
		}
		wg.Wait()

		// assert
		if time.Since(now) > 120 * time.Millisecond {
			t.Error("execution time exceeded")
		}
	})
}
