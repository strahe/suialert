package processors

import (
	"context"
	"fmt"
	"sync"

	"github.com/strahe/suialert/client"
	"github.com/strahe/suialert/config"
	"github.com/strahe/suialert/handlers"
	"github.com/strahe/suialert/types"
	"go.uber.org/zap"
)

type Processor struct {
	cfg *config.Config
	lk  sync.Mutex

	rpcClient *client.Client
	hd        *handlers.SubHandler

	subIDs map[types.EventType]uint64
	done   chan struct{}
}

// NewProcessor creates a new processor
func NewProcessor(cfg *config.Config, rpcClient *client.Client, hd *handlers.SubHandler) (*Processor, error) {
	p := &Processor{
		cfg:       cfg,
		rpcClient: rpcClient,
		hd:        hd,
		subIDs:    make(map[types.EventType]uint64),
		done:      make(chan struct{}),
	}
	return p, nil
}

func (p *Processor) Run(ctx context.Context) error {
	zap.S().Infof("starting processor")
	return p.subscribeEvents(ctx)
}

func (p *Processor) Stop() error {
	close(p.done)
	return p.unsubscribeEvents(context.TODO())
}

func (p *Processor) subscribeEvents(ctx context.Context) error {
	for _, event := range p.cfg.Subscribe.EventTypes {
		zap.S().Infof("subscribing to event: %s", event)
		if err := p.subscribeEventType(ctx, types.EventType(event)); err != nil {
			zap.S().Errorf("failed to subscribe event type %s: %v", event, err)
			return err
		}
	}
	return nil
}

func (p *Processor) unsubscribeEvents(ctx context.Context) error {
	p.lk.Lock()
	defer p.lk.Unlock()

	for event, id := range p.subIDs {
		zap.S().Infof("unsubscribing from event: %s, %d", event, id)
		if ok, err := p.rpcClient.UnsubscribeEvent(ctx, id); err != nil {
			zap.S().Errorf("failed to unsubscribe event type %s: %s", event, err)
			return err
		} else if !ok {
			zap.S().Errorf("failed to unsubscribe event type %s, %d", event, id)
		} else {
			delete(p.subIDs, event)
			zap.S().Infof("unsubscribed from event: %s, %d", event, id)
		}
	}
	return nil
}

// subscribeEventType subscribes to one event type
// https://docs.sui.io/build/event_api#event-filters
func (p *Processor) subscribeEventType(ctx context.Context, eventType types.EventType) error {
	p.lk.Lock()
	defer p.lk.Unlock()

	if sid, ok := p.subIDs[eventType]; ok {
		zap.S().Infof("already subscribed to event type %s: %d", eventType.Name(), sid)
		return nil
	}

	q := types.SubscribeEventQuery{
		EventType: eventType,
	}
	sid, err := p.rpcClient.SubscribeEvent(ctx, q)
	if err != nil {
		return fmt.Errorf("failed to subscribe %s: %s", eventType.Name(), err)
	}

	p.subIDs[eventType] = sid
	zap.S().Infof("subscribed to %s event: %d", eventType.Name(), sid)

	switch eventType {
	case types.EventTypeCoinBalanceChange:
		err = p.hd.AddSub(ctx, eventType.Name(), types.SubscriptionID(sid), p.hd.HandleBalanceChange)
	case types.EventTypePublish:
		err = p.hd.AddSub(ctx, eventType.Name(), types.SubscriptionID(sid), p.hd.HandlePublish)
	case types.EventTypeMoveEvent:
		err = p.hd.AddSub(ctx, eventType.Name(), types.SubscriptionID(sid), p.hd.HandleMove)
	case types.EventTypeNewObject:
		err = p.hd.AddSub(ctx, eventType.Name(), types.SubscriptionID(sid), p.hd.HandleNewObject)
	case types.EventTypeMutateObject:
		err = p.hd.AddSub(ctx, eventType.Name(), types.SubscriptionID(sid), p.hd.HandleMutateObject)
	case types.EventTypeDeleteObject:
		err = p.hd.AddSub(ctx, eventType.Name(), types.SubscriptionID(sid), p.hd.HandleDeleteObject)
	case types.EventTypeTransferObject:
		err = p.hd.AddSub(ctx, eventType.Name(), types.SubscriptionID(sid), p.hd.HandleTransferObject)
	default:
		err = fmt.Errorf("no handler for event: %s", eventType.Name())
	}
	return err
}
