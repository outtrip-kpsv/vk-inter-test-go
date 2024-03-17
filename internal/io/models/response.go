package models

type ErrorResponse struct {
	Error string `json:"error"`
}

type TokenResponse struct {
	Bearer string `json:"bearer"`
}

type OkResponse struct {
	Ok string `json:"ok"`
}
