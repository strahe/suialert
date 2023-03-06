package model

type Rule struct {
	ID         int        `json:"id" pg:"id"`
	Address    string     `json:"address" pg:"address,pk,notnull"`
	UserID     int64      `json:"user_id" pg:"user_id,pk,notnull"`
	User       *User      `json:"user,omitempty" pg:"rel:has-one"`
	AlertLevel AlertLevel `json:"alert_level" pg:"alert_level,notnull"`
	CreatedAt  int64      `json:"created_at" pg:"created_at,notnull"`
	UpdatedAt  int64      `json:"updated_at" pg:"updated_at,notnull"`
}
