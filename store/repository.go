package store

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/floorstyle/fifo-queue/model"
	"github.com/floorstyle/fifo-queue/util"
)

type Repository struct {
	*sync.Mutex
	Cache      model.RedisClient
	Queue      model.FifoQueue
	isConsumed bool
}

const (
	keyList            = "list"
	consumptionTimeout = 5 * time.Second
)

func NewRepository(cache model.RedisClient, queue model.FifoQueue) *Repository {
	return &Repository{
		new(sync.Mutex),
		cache,
		queue,
		false,
	}
}

func (repo *Repository) Init(ctx context.Context) (err error) {
	repo.Lock()
	defer repo.Unlock()

	var msgs []interface{}
	err = repo.Cache.GetMsgs(ctx, keyList, &msgs)
	if err != nil {
		return err
	}
	for _, msg := range msgs {
		err = repo.Queue.PushBackMsg(ctx, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *Repository) GetMessage(ctx context.Context, out interface{}) (err error) {
	repo.Lock()
	defer repo.Unlock()

	if repo.isConsumed {
		return errors.New("already consumed msg that is not deleted")
	}
	err = repo.Cache.GetBackMsg(ctx, keyList, out)
	if err == nil {
		repo.SetAndRefreshConsume(ctx)
		return nil
	}
	util.LogRedisError(err)
	err = repo.Queue.GetBackMsg(ctx, out)
	if err == nil {
		repo.SetAndRefreshConsume(ctx)
	}
	return err
}

func (repo *Repository) PushMessage(ctx context.Context, input interface{}) (err error) {
	repo.Lock()
	defer repo.Unlock()

	util.LogRedisError(repo.Cache.PushFrontMsg(ctx, keyList, input))
	return repo.Queue.PushFrontMsg(ctx, input)
}

func (repo *Repository) PopBackMessage(ctx context.Context, out interface{}) (err error) {
	repo.Lock()
	defer repo.Unlock()

	util.LogRedisError(repo.Cache.PopBackMessage(ctx, keyList, out))
	err = repo.Queue.PopBackMsg(ctx, out)
	if err == nil && repo.isConsumed {
		repo.isConsumed = false
	}
	return err
}

func (repo *Repository) RunAfterCacheIsRestored(ctx context.Context) (err error) {
	for {
		<-repo.Cache.Restored()
		util.LogErr(repo.StoreToCache(ctx))
	}
}

func (repo *Repository) StoreToCache(ctx context.Context) (err error) {
	repo.Lock()
	defer repo.Unlock()

	err = repo.Cache.Delete(ctx, keyList)
	if err != nil {
		return err
	}
	var out []interface{}
	err = repo.Queue.GetMsgs(ctx, &out)
	if err != nil {
		return err
	}
	for _, msg := range out {
		err = repo.Cache.PushFrontMsg(ctx, keyList, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *Repository) SetAndRefreshConsume(ctx context.Context) {
	repo.isConsumed = true

	go util.RunWithBackGround(func(ctx context.Context) error {
		time.Sleep(consumptionTimeout)
		repo.RefreshConsume()
		return nil
	})
}

func (repo *Repository) RefreshConsume() {
	repo.Lock()
	defer repo.Unlock()

	if repo.isConsumed {
		repo.isConsumed = false
	}
}
