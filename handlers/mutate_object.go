package handlers

import (
	"context"

	"github.com/strahe/suialert/model"

	"github.com/strahe/suialert/types"
)

// HandleMutateObject handle mutate object events
func (e *SubHandler) HandleMutateObject(_ context.Context, er *types.EventResult, ed interface{}) error {
	if event, ok := ed.(*types.MutateObject); !ok {
		return nil
	} else {
		if err := e.storeMutateObjectEvent(er, event); err != nil {
			return err
		}
	}
	return nil
}

func (e *SubHandler) storeMutateObjectEvent(er *types.EventResult, ed *types.MutateObject) error {
	m := model.MutateObjectEvent{
		TransactionDigest: er.Id.TxDigest,
		EventSeq:          er.Id.EventSeq,
		Timestamp:         er.Timestamp,
		PackageID:         ed.PackageID,
		TransactionModule: ed.TransactionModule,
		Sender:            types.HexToAddress(ed.Sender),
		ObjectID:          ed.ObjectID,
		ObjectType:        ed.ObjectType,
		Version:           ed.Version,
	}
	return e.db.Create(&m).Error
}
