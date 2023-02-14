package store

import (
	"context"
	"encoding/json"

	"github.com/alicebob/miniredis/v2"
	"github.com/floorstyle/fifo-queue/config"
	"github.com/floorstyle/fifo-queue/model"
	"github.com/floorstyle/fifo-queue/util"
	"github.com/go-redis/redis/v8"
)

type RedisClientReal struct {
	redis.Cmdable
	model.HealthCheck
}

func NewRedisClient(ctx context.Context, config *config.Configuration, out *model.RedisClient) error {
	redisClient := NewRedisConn(config)
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	if out != nil && redisClient != nil {
		*out = redisClient
	}
	return nil
}

func NewRedisConn(config *config.Configuration) *RedisClientReal {
	return &RedisClientReal{
		redis.NewClient(&redis.Options{
			Addr: config.RedisHost,
			DB:   config.RedisDB,
		}),
		NewHealthCheck(config),
	}
}

func (client RedisClientReal) Delete(ctx context.Context, key string) error {
	return client.runWithHealthCheck(ctx, func(ctx context.Context) error {
		return util.AddStack(client.Del(ctx, key).Err())
	})
}

func (client *RedisClientReal) PingClient(ctx context.Context) error {
	err := client.Ping(ctx).Err()
	if err != nil {
		return util.AddStack(err)
	}
	client.HealthCheck.Refresh()
	return nil
}

func (client *RedisClientReal) PushFrontMsg(ctx context.Context, key string, input interface{}) error {
	return client.runWithHealthCheck(ctx, func(ctx context.Context) error {
		value, err := json.Marshal(input)
		if err != nil {
			return util.AddStack(err)
		}
		return util.AddStack(client.LPush(ctx, key, value).Err())
	})
}

func (client *RedisClientReal) GetBackMsg(ctx context.Context, key string, out interface{}) error {
	return client.runWithHealthCheck(ctx, func(ctx context.Context) error {
		result, err := client.LIndex(ctx, key, -1).Bytes()
		if err != nil {
			return err
		}
		return util.AddStack(json.Unmarshal(result, &out))
	})
}

func (client *RedisClientReal) PopBackMessage(ctx context.Context, key string, out interface{}) error {
	return client.runWithHealthCheck(ctx, func(ctx context.Context) error {
		result, err := client.RPop(ctx, key).Bytes()
		if err != nil {
			return err
		}
		return util.AddStack(json.Unmarshal(result, &out))
	})
}

func (client *RedisClientReal) GetMsgs(ctx context.Context, key string, out interface{}) error {
	return client.runWithHealthCheck(ctx, func(ctx context.Context) error {
		msgs, err := client.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			return util.AddStack(err)
		}
		var results []interface{}
		for _, msg := range msgs {
			var res interface{}
			err = json.Unmarshal([]byte(msg), &res)
			if err != nil {
				return util.AddStack(err)
			}

			results = append(results, res)
		}
		return util.UnMarshalStruct(results, out)
	})
}

func (client *RedisClientReal) runWithHealthCheck(ctx context.Context, worker func(context.Context) error) error {
	err := client.HealthCheck.Run(ctx, worker)
	if util.IsErrConnection(err) && !client.HealthCheck.IsHealthy() {
		go util.RunWithBackGroundLoop(util.BACKGROUND_HEALTHCHECK_TIMEOUT, func(ctx context.Context) error {
			err = client.PingClient(ctx)
			if err != nil {
				return err
			}
			client.HealthCheck.RestoreChannel()
			return nil
		})
	}
	return err
}

func NewMockRedisClient(config *config.Configuration, out *model.RedisClient) error {
	mr, err := miniredis.Run()
	if err != nil {
		return err
	}

	redisClient := NewMockRedisConn(mr, config)
	if out != nil && redisClient != nil {
		*out = redisClient
	}
	return nil
}

func NewMockRedisConn(mr *miniredis.Miniredis, config *config.Configuration) *RedisClientReal {
	return &RedisClientReal{
		redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		}),
		NewHealthCheck(config),
	}
}
