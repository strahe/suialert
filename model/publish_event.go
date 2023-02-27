package model

type PublishEvent struct {
	ID EventID `json:"id" pg:"event_id, notnull"`

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
