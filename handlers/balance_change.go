package handlers

import (
	"context"
	"reflect"
	"time"

	"github.com/pgcontrib/bigint"
	"github.com/strahe/suialert/model"
	"github.com/strahe/suialert/types"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// HandleBalanceChange handle the balance change events.
func (e *SubHandler) HandleBalanceChange(ctx context.Context, er *types.EventResult, ed interface{}) error {
	start := time.Now()
	defer func() {
		zap.L().Debug("HandleBalanceChange",
			zap.Duration("took", time.Since(start)),
		)
	}()
	if event, ok := ed.(*types.CoinBalanceChange); !ok {
		zap.S().Errorf("invalid coin balance change event type: %s", reflect.TypeOf(ed))
		return nil
	} else {
		grp, ctx := errgroup.WithContext(ctx)
		grp.Go(func() error {
			return e.eng.ExecuteCoinBalanceChange(ctx, event)
		})
		grp.Go(func() error {
			return e.storeBalanceChangeEvent(ctx, er, event)
		})
		return grp.Wait()
	}
}

func (e *SubHandler) storeBalanceChangeEvent(_ context.Context, er *types.EventResult, ed *types.CoinBalanceChange) error {
	if er == nil || ed == nil {
		return nil
	}
	m := model.CoinBalanceChangeEvent{
		TransactionDigest: er.Id.TxDigest,
		EventSeq:          er.Id.EventSeq,
		Timestamp:         er.Timestamp,
		PackageID:         ed.PackageId,
		TransactionModule: ed.TransactionModule,
		Sender:            types.HexToAddress(ed.Sender),
		ChangeType:        types.CoinBalanceChangeType(ed.ChangeType),
		Owner:             *ed.Owner,
		CoinType:          ed.CoinType,
		CoinObjectID:      ed.CoinObjectId,
		Version:           ed.Version,
		Amount:            bigint.FromInt64(ed.Amount),
	}

	return e.db.Create(&m).Error
}
