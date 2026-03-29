package platformhttp

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeJSONBody(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return ErrInvalidRequestBody
		}

		return err
	}

	return nil
}
