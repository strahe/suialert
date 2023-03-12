package model

import (
	"github.com/strahe/suialert/types"
)

type DeleteObjectEvent struct {
	TransactionDigest string        `json:"tx_digest" gorm:"primaryKey,priority:2"`
	EventSeq          int64         `json:"event_seq"  gorm:"primaryKey,priority:1"`
	Timestamp         uint64        `json:"timestamp"`
	PackageID         string        `json:"package_id"`
	TransactionModule string        `json:"transaction_module"`
	Sender            types.Address `json:"sender"`
	ObjectID          string        `json:"object_id"`
	Version           int64         `json:"version"`
}

func (*DeleteObjectEvent) TableName() string {
	return "delete_object_events"
}
