package schema

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/go-pg/migrations/v8"
)

type patch struct {
	seq  int
	tmpl *template.Template
}

type PatchList struct {
	pm map[int]patch
}

func GetPatches(cfg Config) (*migrations.Collection, error) {
	return patches.Collection(cfg)
}

var patches = NewPatchList()

func NewPatchList() PatchList { // nolint: revive
	return PatchList{map[int]patch{}}
}

// Register adds a patch to the patch list. This should be called in an init function.
func (pl *PatchList) Register(seq int, text string) {
	if seq <= 0 {
		panic(fmt.Sprintf("invalid patch number: %d", seq))
	}

	if _, exists := pl.pm[seq]; exists {
		panic(fmt.Sprintf("duplicate patch registered: %d", seq))
	}

	tmpl, err := template.New("patch").Funcs(schemaTemplateFuncMap).Parse(text)
	if err != nil {
		panic(fmt.Sprintf("parse patch template: %v", err))
	}

	pl.pm[seq] = patch{
		seq:  seq,
		tmpl: tmpl,
	}
}

func (pl *PatchList) Collection(cfg Config) (*migrations.Collection, error) {
	// Check patch list is consistent with no gaps
	count := len(pl.pm)

	// patch 0 must not exist - it's the base schema by definition
	if _, exists := pl.pm[0]; exists {
		return nil, fmt.Errorf("found patch 0, which should not exist")
	}

	// index from 1 since schema seq 0 is the base and not in `pm`
	for i := 1; i <= count; i++ {
		if _, exists := pl.pm[i]; !exists {
			return nil, fmt.Errorf("missing patch %d", i)
		}
	}

	migs := make([]*migrations.Migration, 0, count)
	for i := 1; i <= count; i++ {
		p := pl.pm[i]

		var buf strings.Builder
		if err := p.tmpl.Execute(&buf, cfg); err != nil {
			return nil, fmt.Errorf("execute patch template: %w", err)
		}
		sql := buf.String()

		migs = append(migs, &migrations.Migration{
			Version: int64(i),
			UpTx:    true,
			Up: func(db migrations.DB) error {
				if _, err := db.Exec(sql); err != nil {
					return err
				}
				return nil
			},
		})
	}

	coll := migrations.NewCollection(migs...)
	coll.SetTableName(cfg.SchemaName + ".pg_migrations")
	return coll, nil
}

var schemaTemplateFuncMap = template.FuncMap{
	"default": func(def interface{}, value interface{}) interface{} {
		if isEmpty(value) {
			return def
		}
		return value
	},
}

func isEmpty(val interface{}) bool {
	v := reflect.ValueOf(val)
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Struct:
		return false
	default:
		return v.IsNil()
	}
}
