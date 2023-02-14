package tests

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/floorstyle/fifo-queue/config"
	"github.com/floorstyle/fifo-queue/model"
	"github.com/floorstyle/fifo-queue/server"
	"github.com/floorstyle/fifo-queue/store"
	"github.com/floorstyle/fifo-queue/util"
)

var MockServer server.Server

type T = testing.T
type TB = testing.TB

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	config := config.NewConfiguration()

	var redisClient model.RedisClient
	util.Try(store.NewMockRedisClient(config, &redisClient))

	queue := store.NewFifoQueue()
	repo := store.NewRepository(redisClient, queue)

	MockServer = server.NewMockApiServer(config, repo)
	MockServer.Init(context.Background(), config)
	os.Exit(m.Run())
}

// Initialize context with 1 minute timeout
func testInit(t TB) context.Context {
	const timeOut = 1 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	t.Cleanup(cancel)
	return ctx
}
