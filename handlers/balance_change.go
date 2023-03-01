package handlers

import (
	"context"
	"reflect"

	"github.com/pgcontrib/bigint"
	"github.com/strahe/suialert/model"
	"github.com/strahe/suialert/types"
	"go.uber.org/zap"
)

// HandleBalanceChange handle the balance change events.
func (e *SubHandler) HandleBalanceChange(ctx context.Context, sid types.SubscriptionID, er *types.EventResult, ed interface{}) error {
	if event, ok := ed.(*types.CoinBalanceChange); !ok {
		zap.S().Errorf("invalid coin balance change event type: %s", reflect.TypeOf(ed))
		return nil
	} else {
		if err := e.storeBalanceChangeEvent(ctx, sid, er, event); err != nil {
			zap.S().Errorf("failed to store %s event: %v", e.eventName(sid), err)
		}
	}
	return nil
}

func (e *SubHandler) storeBalanceChangeEvent(_ context.Context, sid types.SubscriptionID, er *types.EventResult, ed *types.CoinBalanceChange) error {
	if er == nil || ed == nil {
		return nil
	}
	m := model.CoinBalanceChangeEvent{
		TransactionDigest: er.Id.TxDigest,
		EventSeq:          er.Id.EventSeq,
		Timestamp:         er.Timestamp,
		PackageID:         ed.PackageId,
		TransactionModule: ed.TransactionModule,
		Sender:            ed.Sender,
		ChangeType:        ed.ChangeType,
		Owner:             types.OwnerToString(ed.Owner),
		CoinType:          ed.CoinType,
		CoinObjectID:      ed.CoinObjectId,
		Version:           ed.Version,
		Amount:            bigint.FromInt64(ed.Amount),
	}
	return e.storeEvent(sid, &m)
}
