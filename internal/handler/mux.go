package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
)

const (
	IncDelta = 1
	DecDelta = -1
)

type muxHandlerRequest struct {
	Urls []string `json:"urls"`
}

type muxHandlerResponse struct {
	Error  bool              `json:"error"`
	Result map[string]string `json:"result"`
}

func (h *Handler) muxHandler(w http.ResponseWriter, r *http.Request) {
	// check count incoming requests
	if h.requestCounter >= h.requestLimit {
		h.errResponse(w,fmt.Errorf("too much requests"), http.StatusInternalServerError)
		return
	}

	// inc & dec request counter
	atomic.AddInt32(&h.requestCounter, IncDelta)
	defer atomic.AddInt32(&h.requestCounter, DecDelta)

	// check method
	if r.Method != http.MethodPost {
		h.errResponse(w,fmt.Errorf("allowed only post method"), http.StatusMethodNotAllowed)
		return
	}

	// parse body
	values := muxHandlerRequest{}
	err := json.NewDecoder(r.Body).Decode(&values)
	if err != nil {
		h.errResponse(w, err, http.StatusBadRequest)
		return
	}

	// validate
	err = h.muxValidate(&values)
	if err != nil {
		h.errResponse(w, err, http.StatusBadRequest)
		return
	}

	// call service
	result, err := h.Mux(r.Context(), values.Urls)
	if err != nil {
		h.errResponse(w, err, http.StatusBadRequest)
		return
	}

	// http response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(muxHandlerResponse{
		Error:  false,
		Result: result,
	})
	if err != nil {
		h.errResponse(w, err, http.StatusInternalServerError)
		return
	}
	return
}
