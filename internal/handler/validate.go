package handler

import (
	"fmt"
	"net/url"
)

const (
	emptyArrayErrorText = "empty array"
	tooMuchUrlErrorText = "urls too much"
	invalidUrlErrorText = "one of the urls is invalid"
)

func (h *Handler) muxValidate(values *muxHandlerRequest) error {
	if len(values.Urls) < 1 {
		return fmt.Errorf(emptyArrayErrorText)
	}

	if int32(len(values.Urls)) > h.urlPerReqLimit {
		return fmt.Errorf(tooMuchUrlErrorText)
	}

	for _, v := range values.Urls {
		u, err := url.Parse(v)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return fmt.Errorf(invalidUrlErrorText)
		}
	}

	return nil
}
