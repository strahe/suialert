package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pgcontrib/bigint"
	"github.com/strahe/suialert/model"
	"github.com/strahe/suialert/types"
	"go.uber.org/zap"
)

// HandleBalanceChange handle the balance change events.
func (e *SubHandler) HandleBalanceChange(ctx context.Context, event *types.Subscription) error {
	if event == nil {
		return nil
	}
	var er *types.EventResult
	if err := json.Unmarshal(event.Result, er); err != nil {
		return fmt.Errorf("error unmarshalling event result: %s", err.Error())
	}

	for name, raw := range er.Event {
		switch name {
		case types.EventCoinBalanceChange.Name():
			var ed *types.CoinBalanceChange
			if err := json.Unmarshal(raw, ed); err != nil {
				zap.S().Errorf("error unmarshalling coin balance change event: %s", err.Error())
				return err
			}
			if err := e.storeBalanceChangeEvent(ctx, er, ed); err != nil {
				zap.S().Errorf("failed to store CoinBalanceChange event: %v", err)
			}
			zap.S().Info(ed)
		default:
			zap.S().Warnf("unknown event name: %s", name)
		}
	}
	return nil
}

func (e *SubHandler) storeBalanceChangeEvent(ctx context.Context, er *types.EventResult, ed *types.CoinBalanceChange) error {
	if er == nil || ed == nil {
		return nil
	}
	cbc := model.CoinBalanceChangeEvent{
		TransactionDigest: er.Id.TxDigest,
		EventSeq:          er.Id.EventSeq,
		Timestamp:         er.Timestamp,
		PackageID:         ed.PackageId,
		TransactionModule: ed.TransactionModule,
		Sender:            ed.Sender,
		ChangeType:        ed.ChangeType,
		Owner:             ed.Owner.AddressOwner,
		CoinType:          ed.CoinType,
		CoinObjectID:      ed.CoinObjectId,
		Version:           ed.Version,
		Amount:            bigint.FromInt64(ed.Amount),
	}
	return e.store.PersistBatch(ctx, &cbc)
}
