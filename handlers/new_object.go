package handlers

import (
	"context"

	"github.com/strahe/suialert/model"

	"go.uber.org/zap"

	"github.com/strahe/suialert/types"
)

// HandleNewObject handle delete object events
func (e *SubHandler) HandleNewObject(ctx context.Context, sid types.SubscriptionID, er *types.EventResult, ed interface{}) error {
	if event, ok := ed.(*types.NewObject); !ok {
		return nil
	} else {
		if err := e.storeNewObjectEvent(ctx, sid, er, event); err != nil {
			zap.S().Errorf("failed to store %s event: %v", e.eventName(sid), err)
		}
	}
	return nil
}

func (e *SubHandler) storeNewObjectEvent(_ context.Context, sid types.SubscriptionID, er *types.EventResult, ed *types.NewObject) error {
	m := model.NewObjectEvent{
		TransactionDigest: er.Id.TxDigest,
		EventSeq:          er.Id.EventSeq,
		Timestamp:         er.Timestamp,
		PackageID:         ed.PackageID,
		TransactionModule: ed.TransactionModule,
		Sender:            ed.Sender,
		Recipient:         types.OwnerToString(ed.Recipient),
		ObjectID:          ed.ObjectID,
		ObjectType:        ed.ObjectType,
		Version:           ed.Version,
	}
	return e.storeEvent(sid, &m)
}
