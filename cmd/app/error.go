package app

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type ServiceError struct {
	Status  int   `json:"status"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func NewServerBadRequestError(message string, errType string) *ServiceError {
	return &ServiceError{Status: fasthttp.StatusBadRequest, Message: message, Type: errType}
}

func NewInternalServerError(message string, errType string) *ServiceError {
	return &ServiceError{Status: fasthttp.StatusInternalServerError, Message: message, Type: errType}
}

func NewServiceError(status int, message string, errType string) *ServiceError {
	return &ServiceError{Status: status, Message: message, Type: errType}
}

func (e *ServiceError) Error() string {
	return e.Message
}

func (e *ServiceError) ToJSON() (string, error) {
	var data []byte
	data, err := json.Marshal(e)

	if err != nil {
		return "", err
	}

	return string(data), nil
}
