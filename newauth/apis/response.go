package apis

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
}

func SetError(w http.ResponseWriter, res *ApiResponse) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	respraw, _ := json.Marshal(&res)
	w.Write(respraw)

}

func SetResponse(code int, w http.ResponseWriter, res interface{}) {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	respraw, _ := json.Marshal(&res)
	w.Write(respraw)

}

func SetSuccessResponse(w http.ResponseWriter) {
	res := ApiResponse{
		Code:    "success",
		Message: "success",
	}

	SetResponse(http.StatusOK, w, res)
}

func SetUserNotFound(w http.ResponseWriter) {
	res := ApiResponse{
		Code:    "user_not_found",
		Message: "user not found",
	}

	SetResponse(http.StatusBadRequest, w, res)
}

func SetValidationErr(w http.ResponseWriter) {
	res := ApiResponse{
		Code:    "validation_error",
		Message: "payload error",
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	respraw, _ := json.Marshal(&res)
	w.Write(respraw)
}
