package model

import (
	"github.com/strahe/suialert/types"
)

type NewObjectEvent struct {
	TransactionDigest string            `json:"tx_digest" gorm:"primaryKey,priority:2"`
	EventSeq          int64             `json:"event_seq"  gorm:"primaryKey,priority:1"`
	Timestamp         uint64            `json:"timestamp"`
	PackageID         string            `json:"package_id"`
	TransactionModule string            `json:"transaction_module"`
	Sender            types.Address     `json:"sender"`
	Recipient         types.ObjectOwner `json:"recipient" gorm:"serializer:json"`
	ObjectType        string            `json:"object_type"`
	ObjectID          string            `json:"object_id"`
	Version           int64             `json:"version"`
}

func (*NewObjectEvent) TableName() string {
	return "new_object_events"
}
