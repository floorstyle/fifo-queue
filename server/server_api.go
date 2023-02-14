package server

import (
	"net/http"

	"github.com/floorstyle/fifo-queue/model"
	"github.com/floorstyle/fifo-queue/util"
	"github.com/gin-gonic/gin"
)

func (this Server) ApiGetMsg(c *gin.Context) {
	ctx := c.Request.Context()

	var response interface{}
	err := this.Repository.GetMessage(ctx, &response)
	if err != nil {
		util.ResponseWithError(c, err, "message")
		return
	}
	c.JSON(http.StatusOK, response)
}

func (this Server) ApiPushMsg(c *gin.Context) {
	ctx := c.Request.Context()

	var input model.Dictionary
	err := util.ValidateBodyInput(c, &input)
	if err != nil {
		util.ResponseWithError(c, err, "message")
		return
	}

	err = this.Repository.PushMessage(ctx, input)
	if err != nil {
		util.ResponseWithError(c, err, "message")
		return
	}
	c.JSON(http.StatusOK, input)
}

func (this Server) ApiRemoveMsg(c *gin.Context) {
	ctx := c.Request.Context()

	var response interface{}
	err := this.Repository.PopBackMessage(ctx, &response)
	if err != nil {
		util.ResponseWithError(c, err, "message")
		return
	}
	c.JSON(http.StatusOK, response)
}
