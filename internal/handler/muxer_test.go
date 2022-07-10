package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_MuxerHandler(t *testing.T) {
	t.Run("wrong method", func(t *testing.T) {
		// arrange
		h := NewHandler()
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
		testServerResponse := muxerHandlerResponse{}
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
		h := NewHandler()
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
		err = json.Unmarshal(testServerResponse.Result[0].Response, &testData)
		if err != nil {
			t.Errorf(err.Error())
		}
		if testData.Value != testContent {
			t.Errorf("expected Value '%s' but got '%s'", testContent, testData.Value)
		}
	})
}
