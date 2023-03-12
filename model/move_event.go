package model

import (
	"github.com/strahe/suialert/types"
)

type MoveEvent struct {
	TransactionDigest string        `json:"tx_digest" gorm:"primaryKey,priority:2"`
	EventSeq          int64         `json:"event_seq"  gorm:"primaryKey,priority:1"`
	Timestamp         uint64        `json:"timestamp"`
	PackageID         string        `json:"package_id"`
	TransactionModule string        `json:"transaction_module"`
	Sender            types.Address `json:"sender"`
	Fields            interface{}   `json:"fields" gorm:"serializer:json"`
	Type              string        `json:"type"`
	BCS               string        `json:"bcs"`
}

func (*MoveEvent) TableName() string {
	return "move_events"
}
