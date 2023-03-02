package storage

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"go.uber.org/zap"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/strahe/suialert/config"
	"github.com/strahe/suialert/model"
	"github.com/strahe/suialert/model/schema"
)

var models = []interface{}{
	(*model.CoinBalanceChangeEvent)(nil),
	(*model.DeleteObjectEvent)(nil),
	(*model.MoveEvent)(nil),
	(*model.MutateObjectEvent)(nil),
	(*model.NewObjectEvent)(nil),
	(*model.PublishEvent)(nil),
	(*model.TransferObjectEvent)(nil),
}

type Database struct {
	db           *pg.DB
	opt          *pg.Options
	schemaConfig schema.Config
	Upsert       bool
}

func NewDatabase(cfg config.PostgresConfig) (*Database, error) {
	opt, err := pg.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}
	if cfg.PoolSize > 0 {
		opt.PoolSize = cfg.PoolSize
	}
	return &Database{
		opt: opt,
		schemaConfig: schema.Config{
			SchemaName: cfg.SchemaName,
		},
		Upsert: cfg.Upsert,
	}, nil
}

func (d *Database) Connect(ctx context.Context) error {
	db, err := connect(ctx, d.opt)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	d.db = db

	return nil
}

func connect(ctx context.Context, opt *pg.Options) (*pg.DB, error) {
	db := pg.Connect(opt)

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

func (d *Database) IsConnected(ctx context.Context) bool {
	if d.db == nil {
		return false
	}

	if err := d.db.Ping(ctx); err != nil {
		return false
	}

	return true
}

func (d *Database) Close() error {
	zap.S().Info("closing database")
	if !d.IsConnected(context.TODO()) {
		err := d.db.Close()
		d.db = nil
		return err
	}
	return nil
}

func (d *Database) SchemaConfig() schema.Config {
	return d.schemaConfig
}

// PersistBatch persists a batch of persistable in a single transaction
func (d *Database) PersistBatch(ctx context.Context, ps ...model.Persistable) error {
	return d.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		txs := &TxStorage{
			tx:     tx,
			upsert: d.Upsert,
		}

		for _, p := range ps {
			if err := p.Persist(ctx, txs); err != nil {
				return fmt.Errorf("persisting %T: %w", p, err)
			}
		}

		return nil
	})
}

type TxStorage struct {
	tx     *pg.Tx
	upsert bool
}

// PersistModel persists a single model
func (s *TxStorage) PersistModel(ctx context.Context, m interface{}) error {
	value := reflect.ValueOf(m)

	elemKind := value.Kind()
	if value.Kind() == reflect.Ptr {
		elemKind = value.Elem().Kind()
	}

	if elemKind == reflect.Slice || elemKind == reflect.Array {
		// Avoid persisting zero length lists
		if value.Len() == 0 {
			return nil
		}

		// go-pg expects pointers to slices. We can fix it up.
		if value.Kind() != reflect.Ptr {
			p := reflect.New(value.Type())
			p.Elem().Set(value)
			m = p.Interface()
		}

	}
	if s.upsert {
		conflict, upsert := GenerateUpsertStrings(m)
		if _, err := s.tx.ModelContext(ctx, m).
			OnConflict(conflict).
			Set(upsert).
			Insert(); err != nil {
			return fmt.Errorf("upserting model: %w", err)
		}
	} else {
		if _, err := s.tx.ModelContext(ctx, m).
			OnConflict("do nothing").
			Insert(); err != nil {
			return fmt.Errorf("persisting model: %w", err)
		}
	}
	return nil
}

func GenerateUpsertStrings(model interface{}) (string, string) {
	var cf []string
	var ucf []string

	// gather all public keys
	for _, pk := range pg.Model(model).TableModel().Table().PKs {
		cf = append(cf, pk.SQLName)
	}
	// gather all other fields
	for _, field := range pg.Model(model).TableModel().Table().DataFields {
		ucf = append(ucf, field.SQLName)
	}

	// consistent ordering in sql statements.
	sort.Strings(cf)
	sort.Strings(ucf)

	// build the conflict string
	var conflict strings.Builder
	conflict.WriteString("(")
	for i, str := range cf {
		conflict.WriteString(str)
		// if this isn't the last field in the conflict statement add a comma.
		if !(i == len(cf)-1) {
			conflict.WriteString(", ")
		}
	}
	conflict.WriteString(") DO UPDATE")

	// build the upsert string
	var upsert strings.Builder
	for i, str := range ucf {
		upsert.WriteString("\"" + str + "\"" + " = EXCLUDED." + str)
		// if this isn't the last field in the upsert statement add a comma.
		if !(i == len(ucf)-1) {
			upsert.WriteString(", ")
		}
	}
	return conflict.String(), upsert.String()
}

func (d *Database) MigrateSchema(ctx context.Context) error {
	db, err := connect(ctx, d.opt)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer db.Close() // nolint: errcheck

	if err := initDatabaseSchema(db, d.SchemaConfig()); err != nil {
		return fmt.Errorf("initializing schema version tables: %w", err)
	}

	return nil
}

// createSchema creates database schema for User and Story models.
func createSchema(db *pg.DB) error {
	for _, mm := range models {
		err := db.Model(mm).CreateTable(&orm.CreateTableOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// initDatabaseSchema initializes the version tables for tracking schema version installed in the database
func initDatabaseSchema(db *pg.DB, cfg schema.Config) error {
	if cfg.SchemaName != "public" {
		_, err := db.Exec(`CREATE SCHEMA IF NOT EXISTS ?`, pg.SafeQuery(cfg.SchemaName))
		if err != nil {
			return fmt.Errorf("ensure schema exists :%w", err)
		}
	}
	return createSchema(db)
}
