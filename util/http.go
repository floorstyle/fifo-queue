package util

import (
	"fmt"
	"net/http"

	"github.com/floorstyle/fifo-queue/model"
	"github.com/gin-gonic/gin"
)

func IsHttpStatusOk(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices
}

func ResponseWithNotFound(c *gin.Context, name string, err error) {
	c.JSON(http.StatusNotFound, gin.H{
		"errorCode": 2048,
		"message":   fmt.Sprintf("The %v cannot be found. Please check input: %v", name, err.Error()),
	})
}

func ResponseWithBadRequest(c *gin.Context, name string, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"errorCode": 2047,
		"message":   fmt.Sprintf("The %v cannot be retrieved due to System Issue: %v", name, err.Error()),
	})
}

func ResponseWithError(c *gin.Context, err error, name string) {
	if IsErrorNotFound(err) {
		ResponseWithNotFound(c, name, err)
	} else {
		ResponseWithBadRequest(c, name, err)
	}
}

func ValidateBodyInput(c *gin.Context, input model.Validator) error {
	err := c.ShouldBindJSON(&input)
	if err != nil {
		return err
	}
	return input.Validate()
}
