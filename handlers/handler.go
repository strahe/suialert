package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/strahe/suialert/bots"
	"github.com/strahe/suialert/model"
	"github.com/strahe/suialert/types"
	client "github.com/strahe/suialert/types"
	"go.uber.org/zap"
)

type handler func(context.Context, types.SubscriptionID, *types.EventResult, interface{}) error

type SubHandler struct {
	handlers   map[types.SubscriptionID]handler
	eventNames map[types.SubscriptionID]string
	lk         sync.Mutex

	bot           bots.Bot
	store         model.Storage
	storeQueued   map[types.SubscriptionID][]model.Persistable
	storeQueuedCh chan *eventPersist
	storeLk       sync.Mutex
	storeBatch    int
	storeLimit    chan struct{}
	done          chan struct{}
}

func NewEthSubHandler(bot bots.Bot, store model.Storage) *SubHandler {
	hd := &SubHandler{
		handlers:      map[client.SubscriptionID]handler{},
		eventNames:    map[client.SubscriptionID]string{},
		bot:           bot,
		store:         store,
		storeQueued:   map[client.SubscriptionID][]model.Persistable{},
		storeQueuedCh: make(chan *eventPersist, 16),
		storeBatch:    100, // todo: make configurable
		storeLimit:    make(chan struct{}, 2),
		done:          make(chan struct{}),
	}
	go hd.background(context.Background())
	return hd
}

func (e *SubHandler) Close() error {
	zap.S().Info("closing subscription handler")
	close(e.done)
	return e.storeAll(context.TODO())
}

func (e *SubHandler) AddSub(name string, id client.SubscriptionID, hd handler) {
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
	return e.eventNames[id]
}

func (e *SubHandler) processSubscription(ctx context.Context, hd handler, p *types.Subscription) error {
	var er types.EventResult
	if err := json.Unmarshal(p.Result, &er); err != nil {
		return fmt.Errorf("error unmarshalling event result: %s", err.Error())
	}

	for name, raw := range er.Event {
		var err error
		switch name {
		case types.EventTypeMutateObject.Name():
			ed := types.MutateObject{}
			if err := json.Unmarshal(raw, &ed); err != nil {
				zap.S().Errorf("error unmarshalling %s event: %s", e.eventName(p.Subscription), err)
				return err
			}
			err = hd(ctx, p.Subscription, &er, &ed)
		case types.EventTypeTransferObject.Name():
			ed := types.TransferObject{}
			if err := json.Unmarshal(raw, &ed); err != nil {
				zap.S().Errorf("error unmarshalling %s event: %s", e.eventName(p.Subscription), err)
				return err
			}
			err = hd(ctx, p.Subscription, &er, &ed)
		case types.EventTypePublish.Name():
			ed := types.Publish{}
			if err := json.Unmarshal(raw, &ed); err != nil {
				zap.S().Errorf("error unmarshalling %s event: %s", e.eventName(p.Subscription), err)
				return err
			}
			err = hd(ctx, p.Subscription, &er, &ed)

		case types.EventTypeCoinBalanceChange.Name():
			ed := types.CoinBalanceChange{}
			if err := json.Unmarshal(raw, &ed); err != nil {
				zap.S().Errorf("error unmarshalling %s event: %s", e.eventName(p.Subscription), err)
				return err
			}
			err = hd(ctx, p.Subscription, &er, &ed)

		case types.EventTypeDeleteObject.Name():
			ed := types.DeleteObject{}
			if err := json.Unmarshal(raw, &ed); err != nil {
				zap.S().Errorf("error unmarshalling %s event: %s", e.eventName(p.Subscription), err)
				return err
			}
			err = hd(ctx, p.Subscription, &er, &ed)
		case types.EventTypeNewObject.Name():
			ed := types.NewObject{}
			if err := json.Unmarshal(raw, &ed); err != nil {
				zap.S().Errorf("error unmarshalling %s event: %s", e.eventName(p.Subscription), err)
				zap.S().Error(string(raw))
				return err
			}
			err = hd(ctx, p.Subscription, &er, &ed)
		case types.EventTypeMoveEvent.Name():
			ed := types.MoveEvent{}
			if err := json.Unmarshal(raw, &ed); err != nil {
				zap.S().Errorf("error unmarshalling %s event: %s", e.eventName(p.Subscription), err)
				return err
			}
			err = hd(ctx, p.Subscription, &er, &ed)
		default:
			zap.S().Warnf("unknown event name %s in %s handler", e.eventName(p.Subscription), name)
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

type eventPersist struct {
	sid types.SubscriptionID
	em  model.Persistable
}

func (e *SubHandler) storeEvent(sid types.SubscriptionID, em model.Persistable) error {
	e.storeQueuedCh <- &eventPersist{sid: sid, em: em}
	return nil
}

func (e *SubHandler) background(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-e.done:
			// clean queue ?
			for {
				select {
				case <-time.After(time.Second * 3):
					return
				case pe := <-e.storeQueuedCh:
					if err := e.appendAndStore(ctx, pe); err != nil {
						zap.S().Errorf("error persisting events: %s", err)
					}
				}
			}

		case pe := <-e.storeQueuedCh:
			if err := e.appendAndStore(ctx, pe); err != nil {
				zap.S().Errorf("error persisting events: %s", err)
			}
		}
	}
}

func (e *SubHandler) appendAndStore(ctx context.Context, pe *eventPersist) error {
	e.storeLk.Lock()
	defer e.storeLk.Unlock()

	e.storeQueued[pe.sid] = append(e.storeQueued[pe.sid], pe.em)

	if len(e.storeQueued[pe.sid]) >= e.storeBatch {
		if err := e.store.PersistBatch(ctx, e.storeQueued[pe.sid]...); err != nil {
			return err
		}
		zap.L().Debug("persisting events",
			zap.String("name", e.eventName(pe.sid)),
			zap.Int("count", len(e.storeQueued[pe.sid])))
		e.storeQueued[pe.sid] = []model.Persistable{}
	}
	return nil
}

func (e *SubHandler) storeAll(ctx context.Context) error {
	e.storeLk.Lock()
	defer e.storeLk.Unlock()

	var err error
	for sid, events := range e.storeQueued {
		if len(events) == 0 {
			continue
		}
		if err := e.store.PersistBatch(ctx, events...); err != nil {
			err = errors.Wrapf(err, "failed to persist %s event to database", e.eventName(sid))
			continue
		}
		zap.L().Debug("persisting events",
			zap.String("name", e.eventName(sid)),
			zap.Int("count", len(events)))
		e.storeQueued[sid] = nil
	}
	return err
}
