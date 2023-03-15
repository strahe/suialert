package model

import (
	"time"

	"github.com/strahe/suialert/types"
	"gorm.io/gorm"
)

type Rule struct {
	ID        uint            `json:"id" gorm:"index"`
	Address   types.Address   `json:"address" gorm:"primaryKey,priority:3"`
	Event     types.EventType `json:"event" gorm:"primaryKey,priority:2"`
	UserID    uint            `json:"user_id" gorm:"primaryKey,autoIncrement:false,priority:1,index"`
	User      User            `json:"-"`
	Condition string          `json:"condition"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (*Rule) TableName() string {
	return "rules"
}

func (r *Rule) AfterCreate(tx *gorm.DB) (err error) {
	return tx.Model(r.User).Update("rule_count", gorm.Expr("rule_count + ?", 1)).Error
}

func (r *Rule) AfterDelete(tx *gorm.DB) (err error) {
	return tx.Model(r.User).Update("rule_count", gorm.Expr("rule_count - ?", 1)).Error
}
