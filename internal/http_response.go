package internal

import "net/http"

func NewHTTPStatus(status int) HTTPStatus {
	return HTTPStatus{
		Status:       status,
		ReasonPhrase: http.StatusText(status),
		Version:      "HTTP/1.1",
	}
}
