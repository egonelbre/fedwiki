package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	Code     int
	Data     interface{}
	Template string
}

func StatusOK(data interface{}) *Response {
	return &Response{Data: data, Template: "", Code: http.StatusOK}
}

type Renderer interface {
	RenderHTML(w io.Writer, template string, data interface{}) error
}

func RenderCommon(rw http.ResponseWriter, r Renderer, responseType string, response *Response) {
	switch {
	case responseType == "" && r == nil:
		responseType = "application/json"
	case responseType == "":
		responseType = "text/html"
	}

	rw.Header().Set("Content-Type", responseType)
	rw.WriteHeader(response.Code)

	switch responseType {
	case "application/json":
		json.NewEncoder(rw).Encode(response.Data)
	case "text/plain":
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(response.Code)
		fmt.Fprintf(rw, "%#v\n", response.Data)
	case "text/html":
		err := r.RenderHTML(rw, response.Template, response.Data)
		if err != nil {
			fmt.Fprintf(rw, err.Error())
		}
	default:
		http.Error(rw, fmt.Sprintf("Unknown Content-Type \"%v\"", responseType), http.StatusNotAcceptable)
	}
}
