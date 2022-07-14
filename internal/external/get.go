package external

import (
	"fmt"
	"io"
	"net/http"
)

func (ex *External) Get(url string) ([]byte, error) {
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	httpResp, err := ex.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.Status != "200 OK" {
		return nil, fmt.Errorf("request with url %s have status code %s", url, httpResp.Status)
	}

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}
