package testhelper

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func CreateReq(method string, path string, payload any) *http.Request {

	data, _ := json.Marshal(&payload)
	req, _ := http.NewRequest(method, path, bytes.NewReader(data))

	return req
}
