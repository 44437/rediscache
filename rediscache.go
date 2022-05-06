package rediscache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
	"time"
)

type Options redis.Options

type cache struct {
	expire time.Duration
	client *redis.Client
}

type Cache interface {
	Get(ctx context.Context, key string, expectedValue interface{}) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}) error
}

func NewRedisCache(
	expire time.Duration,
	redisOptions Options) Cache {

	options := func(redisOptions Options) *redis.Options {
		return &redis.Options{
			Network:            redisOptions.Network,
			Addr:               redisOptions.Addr,
			Dialer:             redisOptions.Dialer,
			OnConnect:          redisOptions.OnConnect,
			Username:           redisOptions.Username,
			Password:           redisOptions.Password,
			DB:                 redisOptions.DB,
			MaxRetries:         redisOptions.MaxRetries,
			MinRetryBackoff:    redisOptions.MinRetryBackoff,
			MaxRetryBackoff:    redisOptions.MaxRetryBackoff,
			DialTimeout:        redisOptions.DialTimeout,
			ReadTimeout:        redisOptions.ReadTimeout,
			WriteTimeout:       redisOptions.WriteTimeout,
			PoolFIFO:           redisOptions.PoolFIFO,
			PoolSize:           redisOptions.PoolSize,
			MinIdleConns:       redisOptions.MinIdleConns,
			MaxConnAge:         redisOptions.MaxConnAge,
			PoolTimeout:        redisOptions.PoolTimeout,
			IdleTimeout:        redisOptions.IdleTimeout,
			IdleCheckFrequency: redisOptions.IdleCheckFrequency,
			TLSConfig:          redisOptions.TLSConfig,
			Limiter:            redisOptions.Limiter,
		}
	}(redisOptions)

	redisClient := redis.NewClient(options)
	return &cache{
		expire: expire,
		client: redisClient,
	}
}

func (c *cache) Get(ctx context.Context, key string, expectedValue interface{}) (
	interface{},
	error,
) {
	var err error
	result, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return expectedValue, err
	}
	var resultMap map[string]interface{}
	err = json.Unmarshal([]byte(result), &resultMap)

	if err != nil {
		return expectedValue, err
	}
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &expectedValue,
		TagName:  "json",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	decoder.Decode(resultMap)

	return expectedValue, err
}

func (c *cache) Set(
	ctx context.Context,
	key string,
	value interface{},
) (err error) {
	valueByte, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.client.Set(ctx, key, valueByte, c.expire*time.Minute)
	return err
}
