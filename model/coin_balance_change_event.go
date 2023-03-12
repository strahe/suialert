package model

import (
	"github.com/pgcontrib/bigint"
	"github.com/strahe/suialert/types"
)

type CoinBalanceChangeEvent struct {
	TransactionDigest string                      `json:"tx_digest" gorm:"primaryKey,priority:2"`
	EventSeq          int64                       `json:"event_seq"  gorm:"primaryKey,priority:1"`
	Timestamp         uint64                      `json:"timestamp"`
	PackageID         string                      `json:"package_id"`
	TransactionModule string                      `json:"transaction_module"`
	Sender            types.Address               `json:"sender"`
	ChangeType        types.CoinBalanceChangeType `json:"change_type"`
	Owner             types.ObjectOwner           `json:"owner" gorm:"serializer:json"`
	CoinType          string                      `json:"coin_type"`
	CoinObjectID      string                      `json:"coin_object_id"`
	Version           int64                       `json:"version"`
	Amount            *bigint.Bigint              `json:"amount"`
}

func (*CoinBalanceChangeEvent) TableName() string {
	return "coin_balance_change_events"
}
