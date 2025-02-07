package server

import "net/http"

const (
	ErrMethodNotAllowed    = "405 Method Not Allowed"
	ErrNotFound            = "404 Not Found"
	ErrInternalServerError = "500 Internal Server Error"
	ErrForbidden           = "403 Forbidden"
)

func sendHTTPError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}
