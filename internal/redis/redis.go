package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	*redis.Client
	ctx context.Context
}

func New(addr, pass string) *Client {
	c := redis.NewClient(&redis.Options{Addr: addr, Password: pass})

	fmt.Printf("Connected to Redis at %s\n", addr)
	return &Client{c, context.Background()}
}

func (c *Client) SetEX(key, val string, ttl time.Duration) error {
	return c.Client.SetEX(c.ctx, key, val, ttl).Err()
}

func (c *Client) Keys(pattern string) ([]string, error) {
	return c.Client.Keys(c.ctx, pattern).Result()
}
