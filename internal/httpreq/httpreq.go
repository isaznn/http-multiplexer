package httpreq

import "net/http"

type HttpReq struct {
	httpClient *http.Client
}

func NewHttpReq(httpClient *http.Client) *HttpReq {
	return &HttpReq{
		httpClient: httpClient,
	}
}
