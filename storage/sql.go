package storage

import (
	"github.com/go-pg/pg/v10"
	"github.com/strahe/suialert/model/schema"
)

type Database struct {
	db           *pg.DB
	opt          *pg.Options
	schemaConfig schema.Config
	Upsert       bool
}
