package helper

import (
	"encoding/json"
	"net/http"
)

// Response struc
type Response struct{}

// StandardResponse is a function for standard response
func (resp Response) StandardResponse(response http.ResponseWriter, code int, v interface{}) {
	var r interface{}
	r = map[string]interface{}{
		"data":            v,
		"responseCode":    code,
		"responseMessage": responseMessage(code),
	}
	json.NewEncoder(response).Encode(r)
}

// StandardResponse is a function for standard response
func (resp Response) StandardResponsePagination(response http.ResponseWriter, code int, v interface{}, p interface{}) {
	var r interface{}
	r = map[string]interface{}{
		"data":            v,
		"pagination":      p,
		"responseCode":    code,
		"responseMessage": responseMessage(code),
	}
	json.NewEncoder(response).Encode(r)
}

// StandardResponseWithMessage is a function for standard response with message
func (resp Response) StandardResponseWithMessage(response http.ResponseWriter, code int, message string) {
	var r interface{}
	r = map[string]interface{}{
		"responseCode":    code,
		"responseMessage": message,
	}
	json.NewEncoder(response).Encode(r)
}

// StandardResponseNoMessage is a function for standard response with message
func (resp Response) StandardResponseNoMessage(response http.ResponseWriter, code int) {
	var r interface{}
	r = map[string]interface{}{
		"responseCode":    code,
		"responseMessage": responseMessage(code),
	}
	json.NewEncoder(response).Encode(r)
}

func responseMessage(code int) string {
	switch code {
	case 0:
		return "Success"
	case 99:
		return "Error"
	default:
		return "Error"
	}
}
