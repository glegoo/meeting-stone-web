package httpserver

import (
	"encoding/json"
	"net/http"
)

type ApiSuccessResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data"`
}

type ApiErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Ok[T any](w http.ResponseWriter, data T) {
	WriteJSON(w, http.StatusOK, ApiSuccessResponse[T]{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func Fail(w http.ResponseWriter, status int, code int, message string, errText string) {
	WriteJSON(w, status, ApiErrorResponse{
		Code:    code,
		Message: message,
		Error:   errText,
	})
}
