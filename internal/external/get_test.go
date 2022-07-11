package external

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
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
		h := NewExternal(&http.Client{})

		// act
		bodyBytes, err := h.Get(testServer.URL)

		// assert
		if err != nil {
			t.Error(err)
		}
		bcomp := bytes.Compare(bodyBytes, []byte(body))
		if bcomp != 0 {
			t.Errorf("expected body '%s' but got '%s'", body, string(bodyBytes))
		}
	})

	t.Run("timeout processing", func(t *testing.T) {
		// arrange
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			time.Sleep(2 * time.Second)
			res.WriteHeader(http.StatusOK)
			res.Write([]byte("ok"))
		}))
		defer func() {
			testServer.Close()
		}()
		h := NewExternal(&http.Client{
			Timeout: 1 * time.Second,
		})

		// act
		_, err := h.Get(testServer.URL)

		// assert
		if !os.IsTimeout(err) {
			t.Error("timeout fail")
		}
	})
}
