package response

import (
	"time"

	"github.com/gin-gonic/gin"
)

type AppResponse struct {
	Status   	int
	Message 	string
	Data			any
	date     	time.Time
}

func NewAppResponse(status int, message string, data any) *AppResponse {
	return &AppResponse{
		Status: status,
		Message: message,
		Data: data,
		date: time.Now(),
	}
}

func (e *AppResponse) Response() map[string]any {
	response := gin.H{
		"status":  e.Status,
		"date":    e.date,
	}

	if e.Data != nil {
		response["data"] = e.Data
	}

	if e.Message != "" {
		response["message"] = e.Message
	}

	return response
}