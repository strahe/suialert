package model

import "context"

type Persistable interface {
	Persist(ctx context.Context, s StorageBatch) error
}

type StorageBatch interface {
	PersistModel(ctx context.Context, m interface{}) error
}

// A Storage can marshal models into a serializable format and persist them.
type Storage interface {
	PersistBatch(ctx context.Context, ps ...Persistable) error
}
