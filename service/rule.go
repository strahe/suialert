package service

import (
	"fmt"

	"github.com/go-pg/pg/v10"

	"github.com/strahe/suialert/model"
)

type RuleService struct {
	db model.Storage
}

func NewRuleService(db model.Storage) *RuleService {
	return &RuleService{db: db}
}

func (s *RuleService) Create(r *model.Rule) error {
	if r == nil {
		return fmt.Errorf("rule is nil")
	}
	_, err := s.db.AsORM().Model(r).Insert()
	return err
}

func (s *RuleService) FindByAddress(uid int64, address string) (*model.Rule, error) {
	var r model.Rule
	err := s.db.AsORM().Model(&r).Where("user_id =?", uid).Where("address =?", address).Select()
	switch err {
	case nil:
		return &r, nil
	case pg.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (s *RuleService) Update(rule *model.Rule) error {
	_, err := s.db.AsORM().Model(rule).WherePK().Column("updated_at", "alert_level").Update()
	return err
}
