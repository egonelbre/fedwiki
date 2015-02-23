package server

import (
	"fmt"
	"net/http"
)

type Error struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Detail string `json:"detail"`
}

func Errorf(code int, format string, args ...interface{}) *Response {
	return &Response{
		Data: Error{
			Status: http.StatusText(code),
			Code:   code,
			Detail: fmt.Sprintf(format, args...),
		},
		Template: "error",
		Code:     code,
	}
}
