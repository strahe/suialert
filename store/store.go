package store

import (
	"context"
	"fmt"
)

var ErrNotFound = fmt.Errorf("not found")

type Store interface {
	Close() error
	Put(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
}
