package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"unicode"
)

const (
	EventTypeMoveEvent         = EventType("MoveEvent")
	EventTypePublish           = EventType("Publish")
	EventTypeCoinBalanceChange = EventType("CoinBalanceChange")
	EventTypeTransferObject    = EventType("TransferObject")
	EventTypeNewObject         = EventType("NewObject")
	EventTypeDeleteObject      = EventType("DeleteObject")
	EventTypeMutateObject      = EventType("MutateObject")
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
	case EventTypeCoinBalanceChange:
		return "Coin balance changing event"
	case EventTypeTransferObject:
		return "Transfer objects to new address / wrap in another object"
	case EventTypeNewObject:
		return "New object creation"
	case EventTypeDeleteObject:
		return "Delete object"
	case EventTypeMutateObject:
		return "Mutate object"
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
	PackageID         string          `json:"package_id"`
	TransactionModule string          `json:"transaction_module"`
	Sender            string          `json:"sender"`
	Type              json.RawMessage `json:"type"`
	Contents          json.RawMessage `json:"contents"`
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
	EventSeq int64  `json:"eventSeq"`
}

// CoinBalanceChange Coin balance changing event
type CoinBalanceChange struct {
	PackageId         string       `json:"packageId"`
	TransactionModule string       `json:"transactionModule"`
	Sender            string       `json:"sender"`
	ChangeType        string       `json:"changeType"`
	Owner             *ObjectOwner `json:"owner"`
	CoinType          string       `json:"coinType"`
	CoinObjectId      string       `json:"coinObjectId"`
	Version           int64        `json:"version"`
	Amount            int64        `json:"amount"`
}

type EventResult struct {
	Timestamp uint64                     `json:"timestamp"`
	TxDigest  string                     `json:"txDigest"`
	Id        EventID                    `json:"id"`
	Event     map[string]json.RawMessage `json:"event"`
}

type ObjectOwner struct {
	*ObjectOwnerInternal
	*string
}

func (o *ObjectOwner) MarshalJSON() ([]byte, error) {
	if o.string != nil {
		data, err := json.Marshal(o.string)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if o.ObjectOwnerInternal != nil {
		data, err := json.Marshal(o.ObjectOwnerInternal)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, errors.New("nil value")
}

func (o *ObjectOwner) UnmarshalJSON(data []byte) error {
	if bytes.HasPrefix(data, []byte("\"")) {
		stringData := string(data[1 : len(data)-1])
		o.string = &stringData
		return nil
	}
	if bytes.HasPrefix(data, []byte("{")) {
		oOI := ObjectOwnerInternal{}
		err := json.Unmarshal(data, &oOI)
		if err != nil {
			return err
		}
		o.ObjectOwnerInternal = &oOI
		return nil
	}
	return errors.New("value not json")
}

func OwnerToString(o *ObjectOwner) string {
	if o == nil {
		return ""
	}
	b, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(b)
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
	PackageID         string       `json:"packageId"`
	TransactionModule string       `json:"transactionModule"`
	Sender            string       `json:"sender"`
	Recipient         *ObjectOwner `json:"recipient"`
	ObjectType        string       `json:"objectType"`
	ObjectID          string       `json:"objectId"`
	Version           int64        `json:"version"`
}

// MutateObject Object level event
// Object mutated.
type MutateObject struct {
	PackageID         string `json:"packageId"`
	TransactionModule string `json:"transactionModule"`
	Sender            string `json:"sender"`
	ObjectType        string `json:"objectType"`
	ObjectID          string `json:"objectId"`
	Version           int64  `json:"version"`
}

// DeleteObject Delete object
type DeleteObject struct {
	PackageID         string `json:"packageId"`
	TransactionModule string `json:"transactionModule"`
	Sender            string `json:"sender"`
	ObjectID          string `json:"objectId"`
	Version           int64  `json:"version"`
}

// NewObject object creation
type NewObject struct {
	PackageID         string       `json:"packageId"`
	TransactionModule string       `json:"transactionModule"`
	Sender            string       `json:"sender"`
	Recipient         *ObjectOwner `json:"recipient"`
	ObjectType        string       `json:"objectType"`
	ObjectID          string       `json:"objectId"`
	Version           int64        `json:"version"`
}
