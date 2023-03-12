package types

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"golang.org/x/crypto/sha3"
)

const (
	AddressLength = 20
)

type ObjectId []byte
type Digest []byte
type Address [AddressLength]byte

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

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) Address { return BytesToAddress(FromHex(s)) }

// String implements fmt.Stringer.
func (a Address) String() string {
	return a.Hex()
}

func (a Address) Hex() string {
	return string(a.checksumHex())
}

func (a Address) hex() []byte {
	var buf [len(a)*2 + 2]byte
	copy(buf[:2], "0x")
	hex.Encode(buf[2:], a[:])
	return buf[:]
}

func (a *Address) checksumHex() []byte {
	buf := a.hex()

	// compute checksum
	sha := sha3.NewLegacyKeccak256()
	sha.Write(buf[2:])
	hash := sha.Sum(nil)
	for i := 2; i < len(buf); i++ {
		hashByte := hash[(i-2)/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if buf[i] > '9' && hashByte > 7 {
			buf[i] -= 32
		}
	}
	return buf[:]
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a), b will be cropped from the left.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *Address) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*a = HexToAddress(str)
	return nil
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (a *Address) Scan(value interface{}) error {
	raw, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var t Address
	err := json.Unmarshal(raw, &t)
	*a = t
	return err
}

// Value return json value, implement driver.Valuer interface
func (a Address) Value() (driver.Value, error) {
	return a.MarshalJSON()
}
