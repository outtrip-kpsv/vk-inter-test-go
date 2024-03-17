package ioutils

import (
	"net/http"
	"vk-inter-test-go/internal/io/models"
)

func HandleInvalidMethodResponse(w http.ResponseWriter, method string) {
	w.WriteHeader(http.StatusBadRequest)
	answer := models.ErrorResponse{
		Error: "Invalid Method: '" + method + "'",
	}
	RespJson(w, answer)
}

func HandleInvalidJson(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	answer := models.ErrorResponse{
		Error: "Error decoding JSON",
	}
	RespJson(w, answer)
}

func RespErrorText(text string, w http.ResponseWriter) {
	answer := models.ErrorResponse{
		Error: text,
	}
	RespJson(w, answer)
}
