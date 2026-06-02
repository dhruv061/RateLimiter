package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse is the standard API response envelope.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// OK sends a successful response with data.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Created sends a 201 response.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Message sends a successful response with a message.
func Message(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: msg,
	})
}

// BadRequest sends a 400 error.
func BadRequest(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Error:   err,
	})
}

// Unauthorized sends a 401 error.
func Unauthorized(c *gin.Context, err string) {
	c.JSON(http.StatusUnauthorized, APIResponse{
		Success: false,
		Error:   err,
	})
}

// Forbidden sends a 403 error.
func Forbidden(c *gin.Context, err string) {
	c.JSON(http.StatusForbidden, APIResponse{
		Success: false,
		Error:   err,
	})
}

// NotFound sends a 404 error.
func NotFound(c *gin.Context, err string) {
	c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Error:   err,
	})
}

// InternalError sends a 500 error.
func InternalError(c *gin.Context, err string) {
	c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Error:   err,
	})
}
