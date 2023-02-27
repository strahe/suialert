package model

type NewObjectEvent struct {
	ID EventID `json:"id" pg:"event_id, notnull"`

	//UTC timestamp in milliseconds
	Timestamp uint64 `json:"timestamp" pg:"timestamp, notnull"`

	// Package ID if available
	PackageID string `json:"package_id" pg:"package_id,"`

	//Module name of the Move package generating the event
	TransactionModule string `json:"transaction_module" pg:"transaction_module,"`

	// Sender in the event
	Sender string `json:"sender" pg:"sender,"`

	//Recipient in the event
	Recipient string `json:"recipient" pg:"recipient,"`

	//Object Type
	ObjectType string `json:"object_type" pg:"object_type,"`

	//Object ID of NewObject, DeleteObject, package being published, or object being transferred
	ObjectID string `json:"object_id" pg:"object_id,"`

	Version int64 `json:"version" pg:"version,"`
}
