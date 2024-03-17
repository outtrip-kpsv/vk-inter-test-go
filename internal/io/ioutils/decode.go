package ioutils

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeRequestBody(req *http.Request, res interface{}) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, res)
	return err
}
