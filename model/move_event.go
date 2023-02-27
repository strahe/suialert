package model

type MoveEvent struct {
	ID EventID `json:"id" pg:"event_id, notnull"`

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
