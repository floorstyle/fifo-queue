package server

import (
	"context"

	"github.com/floorstyle/fifo-queue/config"
	"github.com/floorstyle/fifo-queue/store"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Configuration *config.Configuration
	Router        *gin.Engine
	Repository    *store.Repository
}

func NewApiServer(config *config.Configuration, repo *store.Repository) *Server {
	if config.GinReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	return &Server{
		Configuration: config,
		Router:        gin.Default(),
		Repository:    repo,
	}
}

func (server Server) Init(ctx context.Context, config *config.Configuration) {
	server.Router.GET("/pop", server.ApiGetMsg)
	server.Router.POST("/push", server.ApiPushMsg)
	server.Router.DELETE("/remove", server.ApiRemoveMsg)
}

func NewMockApiServer(config *config.Configuration, repo *store.Repository) Server {
	gin.SetMode(gin.TestMode)
	return Server{
		Configuration: config,
		Router:        gin.Default(),
		Repository:    repo,
	}
}
