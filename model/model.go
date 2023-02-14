package model

import (
	"context"

	"github.com/pkg/errors"
)

type Dictionary map[string]interface{}

func (this Dictionary) IsValid() bool { return len(this) != 0 }

func (this Dictionary) Validate() error {
	if !this.IsValid() {
		return errors.New("input is empty")
	}
	return nil
}

type RedisClient interface {
	PingClient(ctx context.Context) error
	Delete(ctx context.Context, key string) error
	PushFrontMsg(ctx context.Context, key string, input interface{}) error
	GetBackMsg(ctx context.Context, key string, out interface{}) error
	PopBackMessage(ctx context.Context, key string, out interface{}) error
	GetMsgs(ctx context.Context, key string, out interface{}) error
	HealthCheck
}

type FifoQueue interface {
	PushFrontMsg(ctx context.Context, input interface{}) error
	PushBackMsg(ctx context.Context, input interface{}) error
	GetBackMsg(ctx context.Context, out interface{}) error
	PopBackMsg(ctx context.Context, out interface{}) error
	GetMsgs(ctx context.Context, out interface{}) error
}

type HealthCheck interface {
	IsHealthy() bool
	Refresh()
	Run(ctx context.Context, worker func(context.Context) error) error
	Restored() chan bool
	RestoreChannel()
}

type Validator interface {
	Validate() error
}
