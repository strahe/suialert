package model

import "context"

type PublishEvent struct {
	TransactionDigest string `json:"tx_digest" pg:"tx_digest,pk,notnull"`
	EventSeq          int64  `json:"event_seq"  pg:"event_seq,pk,notnull"`

	//UTC timestamp in milliseconds
	Timestamp uint64 `json:"timestamp" pg:"timestamp, notnull"`

	// Sender in the event
	Sender string `json:"sender" pg:"sender,"`

	// Package ID if available
	PackageID string `json:"package_id" pg:"package_id,"`

	//
	Version int64 `json:"version" pg:"version,"`

	//
	Digest string `json:"digest" pg:"digest,"`
}

func (e *PublishEvent) Persist(ctx context.Context, s StorageBatch) error {
	return s.PersistModel(ctx, e)
}
