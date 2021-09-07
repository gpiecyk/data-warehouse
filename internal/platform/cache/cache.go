package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var ErrCacheMiss = errors.New("redis: cache miss")

type Config struct {
	Host string
	Port string

	PoolSize     int
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DialTimeout  time.Duration
}

func NewService(config *Config) (*Client, error) {
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)

	redisPool := &redis.Pool{
		MaxIdle:         config.PoolSize,
		MaxActive:       config.PoolSize,
		IdleTimeout:     config.IdleTimeout,
		Wait:            true,
		MaxConnLifetime: config.IdleTimeout * 2,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				address,
				redis.DialReadTimeout(config.ReadTimeout),
				redis.DialWriteTimeout(config.WriteTimeout),
				redis.DialConnectTimeout(config.DialTimeout),
			)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	return &Client{
		pool: redisPool,
	}, nil
}

type Client struct {
	pool *redis.Pool
}

type Item struct {
	Key   string
	Value interface{}
}

func (c *Client) Set(ctx context.Context, item *Item) error {
	if item.Value == nil {
		return errors.New("nil item value")
	}

	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	value, err := json.Marshal(item.Value)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", item.Key, value)
	return err
}

func (c *Client) Get(ctx context.Context, key string, destination interface{}) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	value, err := redis.Bytes(conn.Do("Get", key))
	if err == redis.ErrNil {
		return ErrCacheMiss
	}
	return json.Unmarshal(value, destination)
}

func (c *Client) Delete(ctx context.Context, key string) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("DEL", key)
	return err
}
