package handlers

import (
	"context"

	"github.com/strahe/suialert/model"

	"github.com/strahe/suialert/types"
)

// HandleNewObject handle delete object events
func (e *SubHandler) HandleNewObject(_ context.Context, er *types.EventResult, ed interface{}) error {
	if event, ok := ed.(*types.NewObject); !ok {
		return nil
	} else {
		if err := e.storeNewObjectEvent(er, event); err != nil {
			return err
		}
	}
	return nil
}

func (e *SubHandler) storeNewObjectEvent(er *types.EventResult, ed *types.NewObject) error {
	m := model.NewObjectEvent{
		TransactionDigest: er.Id.TxDigest,
		EventSeq:          er.Id.EventSeq,
		Timestamp:         er.Timestamp,
		PackageID:         ed.PackageID,
		TransactionModule: ed.TransactionModule,
		Sender:            types.HexToAddress(ed.Sender),
		Recipient:         *ed.Recipient,
		ObjectID:          ed.ObjectID,
		ObjectType:        ed.ObjectType,
		Version:           ed.Version,
	}
	return e.db.Create(&m).Error
}
