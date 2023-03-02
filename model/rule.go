package model

import "github.com/pgcontrib/bigint"

type Role struct {
	Address string      `json:"address" pg:"address,pk,notnull"`
	UserID  int64       `json:"user_id" pg:"user_id,pk,notnull"`
	User    *User       `json:"user,omitempty" pg:"rel:has-one"`
	Rules   []AlertRule `json:"rules" pg:"rules"`
}

type CoinBalanceChangeParams struct {
	ChangeType        *string        `json:"change_type" pg:"change_type,"`
	TransactionModule *string        `json:"transaction_module" pg:"transaction_module,"`
	CoinType          *string        `json:"coin_type" pg:"coin_type,"`
	Amount            *bigint.Bigint `json:"amount" pg:"amount,"`
}

type AlertRule struct {
	Event  *string     `json:"event" pg:"event,notnull"`
	Params interface{} `json:"params" pg:"params"`
}
