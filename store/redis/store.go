package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/strahe/suialert/store"
)

// Client is a redis store client.
type Client struct {
	rdb *redis.Client
}

// NewClient returns a new redis store client.
func NewClient(ctx context.Context, url string) (*Client, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	c := &Client{
		rdb: redis.NewClient(opt),
	}

	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) Close() error {
	if c.rdb != nil {
		return c.rdb.Close()
	}
	return nil
}

// Put stores a key-value pair in the store.
func (c *Client) Put(ctx context.Context, key string, value []byte) error {
	return c.rdb.Set(ctx, key, value, 0).Err()
}

// Get returns the value for the given key.
func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	v, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, store.ErrNotFound
		}
	}
	return v, err
}
