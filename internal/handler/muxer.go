package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/isaznn/http-multiplexer/internal/httpreq"
)

type apiError struct {
	Error        bool   `json:"error"`
	ErrorMessage string `json:"errorMessage"`
}

type urlWithResponse struct {
	Url      string          `json:"url"`
	Response json.RawMessage `json:"response"`
}

type muxerHandlerRequest struct {
	Urls []string `json:"urls"`
}

type muxerHandlerResponse struct {
	Error  bool              `json:"error"`
	Result []urlWithResponse `json:"result"`
}

func (h *Handler) writeErrToJson(w http.ResponseWriter, err error, httpStatus int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(apiError{
		Error:        true,
		ErrorMessage: err.Error(),
	})
}

func (h *Handler) MuxerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,"allowed only post method", http.StatusMethodNotAllowed)
		return
	}

	// parse body
	values := muxerHandlerRequest{}
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

	// TODO move business logic out of package
	resp := muxerHandlerResponse{
		Error: false,
	}
	hr := httpreq.NewHttpReq(&http.Client{
		Timeout: 1 * time.Second,
	})
	for _, v := range values.Urls {
		bodyBytes, err := hr.Get(v)
		if err != nil {
			h.writeErrToJson(w, err, http.StatusBadRequest)
			return
		}
		resp.Result = append(resp.Result, urlWithResponse{
			Url:      v,
			Response: bodyBytes,
		})
	}

	// http response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.writeErrToJson(w, err, http.StatusInternalServerError)
		return
	}
	return
}
