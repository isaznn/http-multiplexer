package external

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

const (
	mockTimeoutMs = 30
)

func TestExternal_Get(t *testing.T) {
	t.Run("broken http client", func(t *testing.T) {
		// arrange
		body := "Hello World"
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(body))
		}))
		defer func() {
			testServer.Close()
		}()
		ex := NewExternal(&http.Client{
			Timeout: mockTimeoutMs * time.Millisecond,
		})

		// act
		bodyBytes, err := ex.Get(testServer.URL)

		// assert
		if err != nil {
			t.Error(err)
		}
		bComp := bytes.Compare(bodyBytes, []byte(body))
		if bComp != 0 {
			t.Errorf("expected body '%s' but got '%s'", body, string(bodyBytes))
		}
	})

	t.Run("timeout processing", func(t *testing.T) {
		// arrange
		testServerSlow := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			time.Sleep(time.Duration(mockTimeoutMs + 5) * time.Millisecond)
			res.WriteHeader(http.StatusOK)
			res.Write([]byte("ok"))
		}))
		defer func() {
			testServerSlow.Close()
		}()
		ex := NewExternal(&http.Client{
			Timeout: mockTimeoutMs * time.Millisecond,
		})

		// act
		_, err := ex.Get(testServerSlow.URL)

		// assert
		if !os.IsTimeout(err) {
			t.Error("timeout fail")
		}
	})

	t.Run("bad status processing", func(t *testing.T) {
		// arrange
		testServerBad := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("error"))
		}))
		defer func() {
			testServerBad.Close()
		}()
		expectedErr := fmt.Errorf("request with url %s have status code 400 Bad Request", testServerBad.URL)
		ex := NewExternal(&http.Client{
			Timeout: mockTimeoutMs * time.Millisecond,
		})

		// act
		_, err := ex.Get(testServerBad.URL)

		// assert
		if !errors.As(err, &expectedErr) {
			t.Error("expected and returned errors do not match")
		}
	})
}
