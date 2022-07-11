package handler

import (
	"encoding/json"
	"net/http"
)

type Muxer interface {
	Mux(urls []string) (map[string]json.RawMessage, error)
}

type Handler struct {
	mux Muxer
}

func NewHandler(m Muxer) *Handler {
	return &Handler{
		mux: m,
	}
}

func (h *Handler) InitRouter() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/muxer", h.muxHandler)
	return r
}
