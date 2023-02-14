package main

import (
	"context"
	"strconv"

	"github.com/floorstyle/fifo-queue/config"
	"github.com/floorstyle/fifo-queue/model"
	"github.com/floorstyle/fifo-queue/server"
	"github.com/floorstyle/fifo-queue/store"
	"github.com/floorstyle/fifo-queue/util"
	"gopkg.in/tylerb/graceful.v1"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := config.NewConfiguration()

	var redisClient model.RedisClient
	util.Try(store.NewRedisClient(ctx, config, &redisClient))

	queue := store.NewFifoQueue()
	repo := store.NewRepository(redisClient, queue)
	util.Try(repo.Init(ctx))
	service := server.NewApiServer(config, repo)
	service.Init(ctx, config)

	go service.Repository.RunAfterCacheIsRestored(ctx)
	graceful.Run(":"+strconv.Itoa(config.HTTPPort), 0, service.Router)
}
