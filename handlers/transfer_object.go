package handlers

import (
	"context"

	"github.com/strahe/suialert/model"

	"github.com/strahe/suialert/types"
)

// HandleTransferObject handle transfer object event
func (e *SubHandler) HandleTransferObject(_ context.Context, er *types.EventResult, ed interface{}) error {
	if event, ok := ed.(*types.TransferObject); !ok {
		return nil
	} else {
		if err := e.storeTransferObjectEvent(er, event); err != nil {
			return err
		}
	}
	return nil
}

func (e *SubHandler) storeTransferObjectEvent(er *types.EventResult, ed *types.TransferObject) error {
	m := model.TransferObjectEvent{
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
