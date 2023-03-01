package model

import "context"

type MoveEvent struct {
	TransactionDigest string `json:"tx_digest" pg:"tx_digest,pk,notnull"`
	EventSeq          int64  `json:"event_seq"  pg:"event_seq,pk,notnull"`

	//UTC timestamp in milliseconds
	Timestamp uint64 `json:"timestamp" pg:"timestamp, notnull"`

	// Package ID if available
	PackageID string `json:"package_id" pg:"package_id,"`

	//Module name of the Move package generating the event
	TransactionModule string `json:"transaction_module" pg:"transaction_module,"`

	// Sender in the event
	Sender string `json:"sender" pg:"sender,"`

	Type []byte `json:"type" pg:"type,"`

	//Contents for MoveEvent
	BCS []byte `json:"bcs" pg:"bcs,"`
}

func (e *MoveEvent) Persist(ctx context.Context, s StorageBatch) error {
	return s.PersistModel(ctx, e)
}
