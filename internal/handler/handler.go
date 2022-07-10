package handler

import "net/http"

type Handler struct {}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) InitRouter() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/muxer", h.MuxerHandler)
	return r
}
