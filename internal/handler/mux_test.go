package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockService struct {}

func (s *MockService) Mux(urls []string) (map[string]json.RawMessage, error) {
	result := make(map[string]json.RawMessage, len(urls))

	for _, v := range urls {
		httpReq, err := http.NewRequest(http.MethodGet, v, nil)
		if err != nil {
			return nil, err
		}

		httpClient := http.Client{Timeout: 1 * time.Second}
		httpResp, err := httpClient.Do(httpReq)
		if err != nil {
			return nil, err
		}
		defer httpResp.Body.Close()

		if httpResp.Status != "200 OK" {
			return nil, fmt.Errorf("request with url %s have status code %s", v, httpResp.Status)
		}

		respBytes, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, err
		}

		result[v] = respBytes
	}

	return result, nil
}


func TestHandler_MuxerHandler(t *testing.T) {
	t.Run("wrong method", func(t *testing.T) {
		// arrange
		s := &MockService{}
		h := NewHandler(s)
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
		s := &MockService{}
		h := NewHandler(s)
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
		err = json.Unmarshal(testServerResponse.Result[testServer.URL], &testData)
		if err != nil {
			t.Errorf(err.Error())
		}
		if testData.Value != testContent {
			t.Errorf("expected Value '%s' but got '%s'", testContent, testData.Value)
		}
	})
}
