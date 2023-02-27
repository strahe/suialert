package model

import "context"

type Persistable interface {
	Persist(ctx context.Context, s StorageBatch) error
}

type StorageBatch interface {
	PersistModel(ctx context.Context, m interface{}) error
}
