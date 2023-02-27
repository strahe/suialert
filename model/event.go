package model

// EventID Unique ID of a Sui Event, the ID is a combination of tx seq number and event seq number,
// the ID is local to this particular fullnode and will be different from other fullnode.
type EventID struct {
	TransactionDigest string `json:"tx_digest" pg:"tx_digest,notnull"`
	EventSeq          int64  `json:"event_seq"  pg:"event_seq,notnull"`
}
