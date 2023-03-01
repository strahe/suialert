package handlers

import (
	"context"

	"github.com/strahe/suialert/model"

	"go.uber.org/zap"

	"github.com/strahe/suialert/types"
)

// HandleMove handle delete object events
func (e *SubHandler) HandleMove(ctx context.Context, sid types.SubscriptionID, er *types.EventResult, ed interface{}) error {
	if event, ok := ed.(*types.DeleteObject); !ok {
		return nil
	} else {
		if err := e.storeDeleteObjectEvent(ctx, sid, er, event); err != nil {
			zap.S().Errorf("failed to store %s event: %v", e.eventName(sid), err)
		}
	}
	return nil
}

func (e *SubHandler) storeMoveEvent(_ context.Context, sid types.SubscriptionID, er *types.EventResult, ed *types.MoveEvent) error {
	m := model.MoveEvent{
		TransactionDigest: er.Id.TxDigest,
		EventSeq:          er.Id.EventSeq,
		Timestamp:         er.Timestamp,
		PackageID:         ed.PackageID,
		TransactionModule: ed.TransactionModule,
		Sender:            ed.Sender,
		Type:              ed.Type,
		BCS:               ed.Contents,
	}
	return e.storeEvent(sid, &m)
}
