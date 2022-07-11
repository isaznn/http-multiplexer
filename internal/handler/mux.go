package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type muxHandlerRequest struct {
	Urls []string `json:"urls"`
}

type muxHandlerResponse struct {
	Error  bool                       `json:"error"`
	Result map[string]json.RawMessage `json:"result"`
}

func (h *Handler) muxHandler(w http.ResponseWriter, r *http.Request) {
	// check method
	if r.Method != http.MethodPost {
		http.Error(w,"allowed only post method", http.StatusMethodNotAllowed)
		return
	}

	// parse body
	values := muxHandlerRequest{}
	err := json.NewDecoder(r.Body).Decode(&values)
	if err != nil {
		h.writeErrToJson(w, err, http.StatusBadRequest)
		return
	}

	// validate
	if len(values.Urls) < 1 {
		h.writeErrToJson(w, fmt.Errorf("empty array"), http.StatusNotAcceptable)
		return
	}
	if len(values.Urls) > 20 {
		h.writeErrToJson(w, fmt.Errorf("urls too much"), http.StatusNotAcceptable)
		return
	}

	// call service
	result, err := h.mux.Mux(values.Urls)
	if err != nil {
		h.writeErrToJson(w, err, http.StatusBadRequest)
		return
	}

	// http response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(muxHandlerResponse{
		Error:  false,
		Result: result,
	})
	if err != nil {
		h.writeErrToJson(w, err, http.StatusInternalServerError)
		return
	}
	return
}
