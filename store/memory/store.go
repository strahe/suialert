package memory

import (
	"context"

	"github.com/strahe/suialert/store"
)

// Client is the memory store client
type Client struct {
	data map[string][]byte
}

// NewClient creates a new memory store client
func NewClient() (*Client, error) {
	return &Client{
		data: map[string][]byte{},
	}, nil
}

func (c *Client) Close() error {
	return nil
}

// Get returns the value for the given key
func (c *Client) Get(_ context.Context, key string) ([]byte, error) {
	if val, ok := c.data[key]; ok {
		return val, nil
	}
	return nil, store.ErrNotFound
}

// Put stores the key/value in the memory store
func (c *Client) Put(_ context.Context, key string, value []byte) error {
	c.data[key] = value
	return nil
}
