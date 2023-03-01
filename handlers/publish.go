package handlers

import (
	"context"

	"github.com/strahe/suialert/model"

	"go.uber.org/zap"

	"github.com/strahe/suialert/types"
)

// HandlePublish handle publish event
func (e *SubHandler) HandlePublish(_ context.Context, sid types.SubscriptionID, er *types.EventResult, ed interface{}) error {
	if event, ok := ed.(*types.Publish); !ok {
		return nil
	} else {
		if err := e.storePublishEvent(sid, er, event); err != nil {
			zap.S().Errorf("failed to store %s event: %v", e.eventName(sid), err)
		}
	}
	return nil
}

func (e *SubHandler) storePublishEvent(sid types.SubscriptionID, er *types.EventResult, ed *types.Publish) error {
	m := model.PublishEvent{
		TransactionDigest: er.Id.TxDigest,
		EventSeq:          er.Id.EventSeq,
		Timestamp:         er.Timestamp,
		PackageID:         ed.PackageID,
		Sender:            ed.Sender,
		Version:           ed.Version,
		Digest:            ed.Digest,
	}
	return e.storeEvent(sid, &m)
}
