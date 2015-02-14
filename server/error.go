package server

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Detail string `json:"detail"`
}

func Error(code int, detail string) (r interface{}, rcode int) {
	return ErrorResponse{
		Status: http.StatusText(code),
		Code:   code,
		Detail: detail,
	}, code
}

func Errorf(code int, format string, args ...interface{}) (r interface{}, rcode int) {
	return ErrorResponse{
		Status: http.StatusText(code),
		Code:   code,
		Detail: fmt.Sprintf(format, args...),
	}, code
}
