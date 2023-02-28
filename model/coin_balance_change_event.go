package model

import (
	"context"
	"github.com/pgcontrib/bigint"
)

//const (
//	CoinBalanceChangeGas CoinBalanceChangeType = iota
//	CoinBalanceChangePay
//	CoinBalanceChangeReceive
//)

type CoinBalanceChangeType int

type CoinBalanceChangeEvent struct {
	TransactionDigest string `json:"tx_digest" pg:"tx_digest,notnull"`
	EventSeq          int64  `json:"event_seq"  pg:"event_seq,notnull"`

	//UTC timestamp in milliseconds
	Timestamp uint64 `json:"timestamp" pg:"timestamp,pk,notnull"`

	// Package ID if available
	PackageID string `json:"package_id" pg:"package_id,"`

	//Module name of the Move package generating the event
	TransactionModule string `json:"transaction_module" pg:"transaction_module,"`

	// Sender in the event
	Sender string `json:"sender" pg:"sender,"`

	ChangeType string `json:"change_type" pg:"change_type,"`

	//Owner in the event
	Owner string `json:"owner" pg:"owner,"`

	CoinType string `json:"coin_type" pg:"coin_type,"`

	CoinObjectID string `json:"coin_object_id" pg:"coin_object_id,"`

	Version int64 `json:"version" pg:"version,"`

	Amount *bigint.Bigint `json:"amount" pg:"amount,"`
}

func (e *CoinBalanceChangeEvent) Persist(ctx context.Context, s StorageBatch) error {
	return s.PersistModel(ctx, e)
}
