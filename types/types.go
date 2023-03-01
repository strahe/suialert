package types

import (
	"encoding/json"
	"errors"
	"reflect"
)

type Address string
type ObjectId []byte
type Digest []byte

type EventQuery struct {
	// Return all events.
	All *string `json:"All"`
	// Return events emitted by the given transaction.
	Transaction *Digest `json:"Transaction"`
	// Return events emitted in a specified Move module
	MoveModule *MoveModule `json:"MoveModule"`
	// Return events with the given move event struct name
	MoveEvent *string `json:"MoveEvent"` // e.g. `0x2::devnet_nft::MintNFTEvent`
	// MoveEvent/Publish/CoinBalanceChange/EpochChange/Checkpoint
	// TransferObject/MutateObject/DeleteObject/NewObject
	EventType *string `json:"EventType"`
	// Query by sender address.
	Sender *Address `json:"Sender"`
	// Query by recipient address
	Recipient *ObjectOwnerInternal `json:"Recipient"`
	// Return events associated with the given object
	Object *ObjectId `json:"Object"`
	// Return events emitted in [start_time, end_time] interval
	TimeRange *TimeRange `json:"TimeRange"`
}

func (q EventQuery) MarshalJSON() ([]byte, error) {
	return marshalQuery(q)
}

func marshalQuery(q any) ([]byte, error) {
	tV := reflect.ValueOf(q)
	for i := 0; i < tV.Type().NumField(); i++ {
		tField := tV.Field(i)
		if tField.Kind() != reflect.Pointer || tField.IsNil() {
			continue
		}
		fieldV := reflect.Indirect(tField)
		tag := tV.Type().Field(i).Tag.Get("json")
		if fieldV.Kind() == reflect.String && tag == "All" {
			return []byte("\"All\""), nil
		}
		data, err := json.Marshal(fieldV.Interface())
		if err != nil {
			return nil, err
		}
		result := []byte("{\"" + tag + "\":")
		result = append(result, data...)
		result = append(result, []byte("}")...)
		return result, nil
	}
	return nil, errors.New("all data is nil")
}

type MoveModule struct {
	Package ObjectId `json:"package"`
	Module  string   `json:"module"`
}

type ObjectOwnerInternal struct {
	AddressOwner *Address `json:"AddressOwner,omitempty"`
	ObjectOwner  *Address `json:"ObjectOwner,omitempty"`
	SingleOwner  *Address `json:"SingleOwner,omitempty"`
	Shared       *struct {
		InitialSharedVersion uint64 `json:"initial_shared_version"`
	} `json:"Shared,omitempty"`
}

type TimeRange struct {
	StartTime uint64 `json:"startTime"` // left endpoint of time interval, milliseconds since epoch, inclusive
	EndTime   uint64 `json:"endTime"`   // right endpoint of time interval, milliseconds since epoch, exclusive
}

type Event json.RawMessage

type EventPage struct {
	Data       []json.RawMessage `json:"data"`
	NextCursor EventID           `json:"nextCursor"`
}

type SubscribeEventQuery struct {
	//Package        *string
	//Module         *string
	//MoveEventType  *string
	//MoveEventField interface{}
	//SenderAddress  *string
	EventType EventType
	//ObjectId       *string
}

type SubscriptionID uint64

type Subscription struct {
	Subscription SubscriptionID  `json:"subscription"`
	Result       json.RawMessage `json:"result"`
}
