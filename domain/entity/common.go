package entity

import (
	"fmt"
	"github.com/vickydk/gosk/utl/config"
	"github.com/vickydk/gosk/utl/log"
	"net/http"
)

type Credentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// response struct
type Response struct {
	Content interface{} `json:"content"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

type ListResponse struct {
	Body interface{} `json:"body"`
	Page interface{} `json:"pagination"`
}

func Respond(err error, result interface{}, code int) *Response {
	if err != nil {
		log.Error(err)
	}
	resp := &Response{
		Content: result,
		Code:    code,
		Message: MappingError(err, code),
	}

	return resp
}

func ListRepond(body interface{}, pagination interface{}) *ListResponse {
	resp := &ListResponse{
		Body: body,
		Page: pagination,
	}

	return resp
}

func MappingError(err error, resCode int) string {
	if config.Env.Debug && err != nil {
		return err.Error()
	} else {
		switch resCode {
		case http.StatusOK:
			return ""
		case http.StatusPaymentRequired:
			return "Payment Required"
		case http.StatusBadRequest:
			return "Missing request"
		case http.StatusConflict:
			return "Data Conflict"
		case http.StatusNotFound:
			return "Data Not Found"
		case http.StatusInsufficientStorage:
			return "Error Database"
		case http.StatusUnauthorized:
			return "User is not authorized"
		default:
			return fmt.Sprint("Not mapping yet: ", resCode)
		}
	}
}
