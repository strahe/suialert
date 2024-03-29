package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/strahe/suialert/types"
	"gorm.io/gorm"

	"github.com/strahe/suialert/model"
)

type RuleService struct {
	db *gorm.DB
}

func NewRuleService(db *gorm.DB) *RuleService {
	return &RuleService{db: db}
}

func (s *RuleService) Create(r *model.Rule) error {
	if r == nil {
		return fmt.Errorf("rule is nil")
	}
	return s.db.Create(r).Error
}

func (s *RuleService) FindByPrimaryKey(uid uint, event types.EventType, addr types.Address) (*model.Rule, error) {
	var rule model.Rule
	err := s.db.Where(&model.Rule{UserID: uid, Event: event, Address: addr},
		"UserID", "Event", "Address").First(&rule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &rule, nil
}

func (s *RuleService) Update(rule *model.Rule) error {
	if rule == nil {
		return fmt.Errorf("rule is nil")
	}
	return s.db.Select("Condition").Updates(rule).Error
}

func (s *RuleService) FindByAddress(addr types.Address) ([]model.Rule, error) {
	var rules []model.Rule
	if err := s.db.Where("address = ?", addr).Find(&rules).Error; err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	return rules, nil
}

// FindAll returns all rules
// todo: pagination
func (s *RuleService) FindAll(ctx context.Context) ([]model.Rule, error) {
	var rules []model.Rule
	if err := s.db.WithContext(ctx).Find(&rules).Error; err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	return rules, nil
}
