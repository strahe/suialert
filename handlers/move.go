package handlers

import (
	"context"

	"github.com/strahe/suialert/model"

	"go.uber.org/zap"

	"github.com/strahe/suialert/types"
)

// HandleMove handle delete object events
func (e *SubHandler) HandleMove(_ context.Context, sid types.SubscriptionID, er *types.EventResult, ed interface{}) error {
	if event, ok := ed.(*types.MoveEvent); !ok {
		return nil
	} else {
		if err := e.storeMoveEvent(er, event); err != nil {
			zap.S().Errorf("failed to store %s event: %v", e.eventName(sid), err)
		}
	}
	return nil
}

func (e *SubHandler) storeMoveEvent(er *types.EventResult, ed *types.MoveEvent) error {
	m := model.MoveEvent{
		TransactionDigest: er.Id.TxDigest,
		EventSeq:          er.Id.EventSeq,
		Timestamp:         er.Timestamp,
		PackageID:         ed.PackageId,
		TransactionModule: ed.TransactionModule,
		Sender:            types.HexToAddress(ed.Sender),
		Fields:            ed.Fields,
		Type:              ed.Type,
		BCS:               ed.Bcs,
	}
	return e.db.Create(&m).Error
}
