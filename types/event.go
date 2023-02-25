package types

import (
	"encoding/json"
	"unicode"
)

const (
	Gas = iota
	Pay
	Receive
)

const (
	EventTypeMoveEvent     = EventType("MoveEvent")
	EventTypePublish       = EventType("Publish")
	EventCoinBalanceChange = EventType("CoinBalanceChange")
	EventTransferObject    = EventType("TransferObject")
	EventNewObject         = EventType("NewObject")
	EventEpochChange       = EventType("EpochChange")
	EventCheckpoint        = EventType("Checkpoint")
)

type EventType string

func (e EventType) Name() string {
	r := []rune(e)
	return string(append([]rune{unicode.ToLower(r[0])}, r[1:]...))
}

func (e EventType) Description() string {
	switch e {
	case EventTypeMoveEvent:
		return "Move-specific event"
	case EventTypePublish:
		return "Module published"
	case EventCoinBalanceChange:
		return "Coin balance changing event"
	case EventTransferObject:
		return "Transfer objects to new address / wrap in another object"
	case EventNewObject:
		return "New object creation"
	case EventEpochChange:
		return "Epoch change"
	case EventCheckpoint:
		return "New checkpoint"
	}
	return "Unknown event"
}

type StructTag struct {
	Address    string `json:"address"`
	Module     string `json:"module"`
	Name       string `json:"name"`
	TypeParams string `json:"type_params"`
}

// MoveEvent is a move event.
// Transaction level event
// Move-specific event
type MoveEvent struct {
	PackageID         string    `json:"package_id"`
	TransactionModule string    `json:"transaction_module"`
	Sender            string    `json:"sender"`
	Type              StructTag `json:"type"`
	Contents          string    `json:"contents"`
}

// Publish Module published
type Publish struct {
	Sender    string `json:"sender"`
	PackageID string `json:"package_id"`
	Version   int64  `json:"version"`
	Digest    string `json:"digest"`
}

type EventID struct {
	TxDigest string `json:"txDigest"`
	EventSeq int    `json:"eventSeq"`
}

// CoinBalanceChange Coin balance changing event
type CoinBalanceChange struct {
	PackageId         string `json:"packageId"`
	TransactionModule string `json:"transactionModule"`
	Sender            string `json:"sender"`
	ChangeType        string `json:"changeType"`
	Owner             Owner  `json:"owner"`
	CoinType          string `json:"coinType"`
	CoinObjectId      string `json:"coinObjectId"`
	Version           int    `json:"version"`
	Amount            int    `json:"amount"`
}

type EventResult struct {
	Timestamp int64                      `json:"timestamp"`
	TxDigest  string                     `json:"txDigest"`
	Id        EventID                    `json:"id"`
	Event     map[string]json.RawMessage `json:"event"`
}

type Owner struct {
	AddressOwner string `json:"AddressOwner"`
}

// EpochChange Epoch change
type EpochChange struct {
	EpochId uint64 `json:"epoch_id"`
}

// Checkpoint New checkpoint
type Checkpoint struct {
	CheckpointSequenceNumber uint64 `json:"checkpoint_sequence_number"`
}

// TransferObject Object level event
// Transfer objects to new address / wrap in another object
type TransferObject struct {
	PackageID         string `json:"package_id"`
	TransactionModule string `json:"transaction_module"`
	Sender            string `json:"sender"`
	Recipient         string `json:"recipient"`
	ObjectType        string `json:"object_type"`
	ObjectID          string `json:"object_id"`
	Version           uint64 `json:"version"`
}

// MutateObject Object level event
// Object mutated.
type MutateObject struct {
	PackageID         string `json:"package_id"`
	TransactionModule string `json:"transaction_module"`
	Sender            string `json:"sender"`
	ObjectType        string `json:"object_type"`
	ObjectID          string `json:"object_id"`
	Version           uint64 `json:"version"`
}

// DeleteObject Delete object
type DeleteObject struct {
	PackageID         string `json:"package_id"`
	TransactionModule string `json:"transaction_module"`
	Sender            string `json:"sender"`
	ObjectID          string `json:"object_id"`
	Version           uint64 `json:"version"`
}

// NewObject object creation
type NewObject struct {
	PackageID         string `json:"package_id"`
	TransactionModule string `json:"transaction_module"`
	Sender            string `json:"sender"`
	Recipient         string `json:"recipient"`
	ObjectType        string `json:"object_type"`
	ObjectID          string `json:"object_id"`
	Version           uint64 `json:"version"`
}
