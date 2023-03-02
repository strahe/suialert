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

type AlertLevel string

const (
	AlertLevelNone   AlertLevel = "None"   // not alerting
	AlertLevelLow    AlertLevel = "Low"    //
	AlertLevelMedium AlertLevel = "Medium" //
	AlertLevelHigh   AlertLevel = "High"   //
	AlertLevelAll    AlertLevel = "All"    // send all alerts
)
