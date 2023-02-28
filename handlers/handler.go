package handlers

import (
	"context"
	"sync"

	"github.com/strahe/suialert/model"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/strahe/suialert/bots"
	"github.com/strahe/suialert/types"
	client "github.com/strahe/suialert/types"
)

type SubHandler struct {
	queued map[types.SubscriptionID][]types.Subscription
	sinks  map[types.SubscriptionID]func(context.Context, *types.Subscription) error
	lk     sync.Mutex

	bot   bots.Bot
	store model.Storage
}

func NewEthSubHandler(bot bots.Bot, store model.Storage) *SubHandler {
	return &SubHandler{
		queued: map[types.SubscriptionID][]types.Subscription{},
		sinks:  map[client.SubscriptionID]func(context.Context, *types.Subscription) error{},
		bot:    bot,
		store:  store,
	}
}

func (e *SubHandler) AddSub(ctx context.Context, id client.SubscriptionID, sink func(context.Context, *types.Subscription) error) error {
	e.lk.Lock()
	defer e.lk.Unlock()

	for _, p := range e.queued[id] {
		p := p // copy
		if err := sink(ctx, &p); err != nil {
			return err
		}
	}
	delete(e.queued, id)
	e.sinks[id] = sink
	return nil
}

func (e *SubHandler) RemoveSub(id client.SubscriptionID) {
	e.lk.Lock()
	defer e.lk.Unlock()

	delete(e.sinks, id)
	delete(e.queued, id)
}

func (e *SubHandler) SubscribeEvent(ctx context.Context, r jsonrpc.RawParams) error {
	p, err := jsonrpc.DecodeParams[client.Subscription](r)
	if err != nil {
		return err
	}

	e.lk.Lock()

	sink := e.sinks[p.Subscription]

	if sink == nil {
		e.queued[p.Subscription] = append(e.queued[p.Subscription], p)
		e.lk.Unlock()
		return nil
	}

	e.lk.Unlock()
	return sink(ctx, &p) // todo track errors and auto-unsubscribe on rpc conn close?
}
