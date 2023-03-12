package model

import (
	"github.com/strahe/suialert/types"
)

type PublishEvent struct {
	TransactionDigest string        `json:"tx_digest" gorm:"primaryKey,priority:2"`
	EventSeq          int64         `json:"event_seq"  gorm:"primaryKey,priority:1"`
	Timestamp         uint64        `json:"timestamp"`
	Sender            types.Address `json:"sender"`
	PackageID         string        `json:"package_id"`
	Version           int64         `json:"version"`
	Digest            string        `json:"digest"`
}

func (*PublishEvent) TableName() string {
	return "publish_events"
}
