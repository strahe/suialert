package model

import (
	"bytes"
	"text/template"
	"time"

	"github.com/strahe/suialert/types"
	"gorm.io/gorm"
)

type Rule struct {
	ID        uint            `json:"id" gorm:"index"`
	Address   types.Address   `json:"address" gorm:"primaryKey,priority:3,index"`
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

var (
	gtp *template.Template
)

func init() {
	grl := `
rule Rule{{ .ID }} "" salience 10  {
    when
        {{ .Condition }}
    then
		Event.CoinType = "xxxx";
}
`
	t, err := template.New("rule").Parse(grl)
	if err != nil {
		panic(err)
	}
	gtp = t
}

func (r *Rule) BuildGRL() ([]byte, error) {
	var buf bytes.Buffer
	if err := gtp.Execute(&buf, r); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
