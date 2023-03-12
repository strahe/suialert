package handlers

import (
	"context"

	"github.com/strahe/suialert/model"

	"go.uber.org/zap"

	"github.com/strahe/suialert/types"
)

// HandleDeleteObject handle delete object events
func (e *SubHandler) HandleDeleteObject(ctx context.Context, sid types.SubscriptionID, er *types.EventResult, ed interface{}) error {
	if event, ok := ed.(*types.DeleteObject); !ok {
		return nil
	} else {
		if err := e.storeDeleteObjectEvent(ctx, er, event); err != nil {
			zap.S().Errorf("failed to store %s event: %v", e.eventName(sid), err)
		}
	}
	return nil
}

func (e *SubHandler) storeDeleteObjectEvent(_ context.Context, er *types.EventResult, ed *types.DeleteObject) error {
	m := model.DeleteObjectEvent{
		TransactionDigest: er.Id.TxDigest,
		EventSeq:          er.Id.EventSeq,
		Timestamp:         er.Timestamp,
		PackageID:         ed.PackageID,
		TransactionModule: ed.TransactionModule,
		Sender:            types.HexToAddress(ed.Sender),
		ObjectID:          ed.ObjectID,
		Version:           ed.Version,
	}
	return e.db.Create(&m).Error
}
