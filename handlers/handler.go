package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/strahe/suialert/rule"
	"gorm.io/gorm"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/strahe/suialert/bots"
	"github.com/strahe/suialert/types"
	client "github.com/strahe/suialert/types"
	"go.uber.org/zap"
)

type handler func(context.Context, *types.EventResult, interface{}) error

type SubHandler struct {
	handlers   map[types.SubscriptionID]handler
	eventNames map[types.SubscriptionID]types.EventType
	lk         sync.Mutex

	bot  bots.Bot
	db   *gorm.DB
	eng  *rule.Engine
	done chan struct{}
}

func NewSubHandler(bot bots.Bot, db *gorm.DB, eng *rule.Engine) *SubHandler {
	hd := &SubHandler{
		handlers:   map[client.SubscriptionID]handler{},
		eventNames: map[client.SubscriptionID]types.EventType{},
		bot:        bot,
		db:         db,
		eng:        eng,
		done:       make(chan struct{}),
	}
	return hd
}

func (e *SubHandler) Close() error {
	zap.S().Info("closing subscription handler")
	close(e.done)
	return nil
}

func (e *SubHandler) AddSub(name types.EventType, id client.SubscriptionID, hd handler) {
	e.lk.Lock()
	defer e.lk.Unlock()

	e.handlers[id] = hd
	e.eventNames[id] = name
}

func (e *SubHandler) RemoveSub(id client.SubscriptionID) {
	e.lk.Lock()
	defer e.lk.Unlock()

	delete(e.handlers, id)
	delete(e.eventNames, id)
}

func (e *SubHandler) SubscribeEvent(ctx context.Context, r jsonrpc.RawParams) error {
	p, err := jsonrpc.DecodeParams[client.Subscription](r)
	if err != nil {
		return err
	}

	e.lk.Lock()

	hd := e.handlers[p.Subscription]
	if hd == nil {
		e.lk.Unlock()
		return nil
	}
	e.lk.Unlock()
	return e.processSubscription(ctx, hd, &p)
}

func (e *SubHandler) eventName(id types.SubscriptionID) string {
	return string(e.eventNames[id])
}

func (e *SubHandler) processSubscription(ctx context.Context, hd handler, p *types.Subscription) error {
	var er types.EventResult
	if err := json.Unmarshal(p.Result, &er); err != nil {
		return fmt.Errorf("error unmarshalling event result: %s", err.Error())
	}

	for name, raw := range er.Event {
		var err error
		switch types.EventFromSui(name) {
		case types.EventTypeMutateObject:
			ed := types.MutateObject{}
			if err := json.Unmarshal(raw, &ed); err != nil {
				zap.S().Errorf("error unmarshalling %s event: %s", e.eventName(p.Subscription), err)
				return err
			}
			err = hd(ctx, &er, &ed)
		case types.EventTypeTransferObject:
			ed := types.TransferObject{}
			if err := json.Unmarshal(raw, &ed); err != nil {
				zap.S().Errorf("error unmarshalling %s event: %s", e.eventName(p.Subscription), err)
				return err
			}
			err = hd(ctx, &er, &ed)
		case types.EventTypePublish:
			ed := types.Publish{}
			if err := json.Unmarshal(raw, &ed); err != nil {
				zap.S().Errorf("error unmarshalling %s event: %s", e.eventName(p.Subscription), err)
				return err
			}
			err = hd(ctx, &er, &ed)

		case types.EventTypeCoinBalanceChange:
			ed := types.CoinBalanceChange{}
			if err := json.Unmarshal(raw, &ed); err != nil {
				zap.S().Errorf("error unmarshalling %s event: %s", e.eventName(p.Subscription), err)
				return err
			}
			err = hd(ctx, &er, &ed)

		case types.EventTypeDeleteObject:
			ed := types.DeleteObject{}
			if err := json.Unmarshal(raw, &ed); err != nil {
				zap.S().Errorf("error unmarshalling %s event: %s", e.eventName(p.Subscription), err)
				return err
			}
			err = hd(ctx, &er, &ed)
		case types.EventTypeNewObject:
			ed := types.NewObject{}
			if err := json.Unmarshal(raw, &ed); err != nil {
				zap.S().Errorf("error unmarshalling %s event: %s", e.eventName(p.Subscription), err)
				return err
			}
			err = hd(ctx, &er, &ed)
		case types.EventTypeMove:
			ed := types.MoveEvent{}
			if err := json.Unmarshal(raw, &ed); err != nil {
				zap.S().Errorf("error unmarshalling %s event: %s", e.eventName(p.Subscription), err)
				return err
			}
			err = hd(ctx, &er, &ed)
		default:
			zap.S().Warnf("unknown event name %s in %s handler", e.eventName(p.Subscription), name)
			continue
		}
		if err != nil {
			zap.L().Error("error processing event",
				zap.String("name", e.eventName(p.Subscription)),
				zap.Error(err))
			return err
		}
	}
	return nil
}
