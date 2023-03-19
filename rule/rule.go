package rule

import (
	"context"
	"fmt"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/samber/lo"
	"github.com/strahe/suialert/service"
	"github.com/strahe/suialert/types"
	"go.uber.org/zap"
)

type Engine struct {
	eg  *engine.GruleEngine
	lib *ast.KnowledgeLibrary

	rsv *service.RuleService
}

func NewEngine(rsv *service.RuleService) (*Engine, error) {
	e := Engine{
		eg:  engine.NewGruleEngine(),
		lib: ast.NewKnowledgeLibrary(),

		rsv: rsv,
	}
	return &e, nil
}

func (e *Engine) LoadRules(ctx context.Context) error {
	rules, err := e.rsv.FindAll(ctx)
	if err != nil {
		return err
	}
	rb := builder.NewRuleBuilder(e.lib)

	for _, rule := range rules {
		kb := e.lib.NewKnowledgeBaseInstance(rule.Address.Hex(), "")
		if kb == nil {
			gpl, err := rule.BuildGRL()
			if err != nil {
				return err
			}
			if err := rb.BuildRuleFromResource(rule.Address.Hex(), "", pkg.NewBytesResource(gpl)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Engine) ExecuteCoinBalanceChange(_ context.Context, data *types.CoinBalanceChange) error {
	dataCtx := ast.NewDataContext()
	err := dataCtx.Add("Event", data)
	if err != nil {
		return err
	}
	owner, ok := lo.Coalesce[*types.Address](data.Owner.ObjectOwner, data.Owner.AddressOwner, data.Owner.SingleOwner)
	if !ok {
		// todo: handle this case
		return fmt.Errorf("invalid owner")
	}
	knowledgeBase := e.lib.NewKnowledgeBaseInstance(owner.Hex(), "")
	if knowledgeBase == nil {
		zap.S().Debugf("no rules matchd for owner: %s", owner.Hex())
		return nil
	}
	rules, err := e.eg.FetchMatchingRules(dataCtx, knowledgeBase)
	if err != nil {
		return err
	}
	fmt.Println("match rules", rules)
	return nil
}
