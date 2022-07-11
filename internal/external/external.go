package external

import "net/http"

type External struct {
	httpClient *http.Client
}

func NewExternal(httpClient *http.Client) *External {
	return &External{
		httpClient: httpClient,
	}
}
