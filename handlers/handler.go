package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/strahe/suialert/bots"
	"github.com/strahe/suialert/model"
	"github.com/strahe/suialert/types"
	client "github.com/strahe/suialert/types"
	"go.uber.org/zap"
)

type handler func(context.Context, types.SubscriptionID, *types.EventResult, interface{}) error

type SubHandler struct {
	queued     map[types.SubscriptionID][]types.Subscription
	handlers   map[types.SubscriptionID]handler
	eventNames map[types.SubscriptionID]string
	lk         sync.Mutex

	bot           bots.Bot
	store         model.Storage
	storeQueued   map[types.SubscriptionID][]model.Persistable
	storeQueuedCh chan *eventPersist
	storeLk       sync.Mutex
	storeLimit    int
}

func NewEthSubHandler(ctx context.Context, bot bots.Bot, store model.Storage) *SubHandler {
	hd := &SubHandler{
		queued:        map[types.SubscriptionID][]types.Subscription{},
		handlers:      map[client.SubscriptionID]handler{},
		eventNames:    map[client.SubscriptionID]string{},
		bot:           bot,
		store:         store,
		storeQueued:   map[client.SubscriptionID][]model.Persistable{},
		storeQueuedCh: make(chan *eventPersist, 16),
		storeLimit:    100, // todo: make configurable
	}
	go hd.doStore(ctx)
	return hd
}

func (e *SubHandler) AddSub(ctx context.Context, name string, id client.SubscriptionID, hd handler) error {
	e.lk.Lock()
	defer e.lk.Unlock()

	for _, p := range e.queued[id] {
		p := p // copy
		if err := e.processSubscription(ctx, hd, &p); err != nil {
			return err
		}
	}
	delete(e.queued, id)
	e.handlers[id] = hd
	e.eventNames[id] = name
	return nil
}

func (e *SubHandler) RemoveSub(id client.SubscriptionID) {
	e.lk.Lock()
	defer e.lk.Unlock()

	delete(e.handlers, id)
	delete(e.queued, id)
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
		e.queued[p.Subscription] = append(e.queued[p.Subscription], p)
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

	//mk := func(ed interface{}, raw json.RawMessage) handler {
	//	if err := json.Unmarshal(raw, &ed); err != nil {
	//		zap.S().Errorf("error unmarshalling %s event: %s", e.eventName(p.Subscription), err)
	//		return func(ctx context.Context, id client.SubscriptionID, result *client.EventResult, i interface{}) error {
	//			return err
	//		}
	//	}
	//	return func(ctx context.Context, id client.SubscriptionID, result *client.EventResult, i interface{}) error {
	//		return hd(ctx, p.Subscription, &er, &ed)
	//	}
	//}

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
	e.storeQueuedCh <- &eventPersist{
		sid: sid,
		em:  em,
	}
	return nil
}

func (e *SubHandler) doStore(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// todo: store all left items
			return
		case eb := <-e.storeQueuedCh:
			e.storeLk.Lock()
			e.storeQueued[eb.sid] = append(e.storeQueued[eb.sid], eb.em)
			if len(e.storeQueued[eb.sid]) >= e.storeLimit {
				events := e.storeQueued[eb.sid]
				zap.S().Debugf("store %d %s events to db", len(events), e.eventName(eb.sid))
				if err := e.store.PersistBatch(ctx, events...); err != nil {
					zap.S().Errorf("failed to store %s event to db: %v", e.eventName(eb.sid), err)
					e.storeLk.Unlock()
					continue
				}
				e.storeQueued[eb.sid] = []model.Persistable{}
			}
			e.storeLk.Unlock()
		}
	}
}
