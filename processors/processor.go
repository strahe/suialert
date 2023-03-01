package processors

import (
	"context"
	"fmt"
	"sync"
	"time"

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

func (p *Processor) Close() error {
	close(p.done)
	return p.unsubscribeEvents(context.TODO())
}

func (p *Processor) subscribeEvents(ctx context.Context) error {
	for _, event := range p.cfg.Subscribe.EventTypes {
		zap.L().Info("subscribing",
			zap.String("event", event),
			zap.Time("start", time.Now()),
		)
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
		zap.L().Info("unsubscribing",
			zap.String("event", event.Name()),
			zap.Uint64("id", id),
		)

		if ok, err := p.rpcClient.UnsubscribeEvent(ctx, id); err != nil {
			zap.S().Errorf("failed to unsubscribe event type %s: %s", event, err)
			return err
		} else if !ok {
			zap.S().Errorf("failed to unsubscribe event type %s, %d", event, id)
		} else {
			delete(p.subIDs, event)
			zap.L().Info("unsubscribed",
				zap.String("event", event.Name()), //
				zap.Uint64("id", id),              //
			)
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
	zap.L().Info("subscribed",
		zap.String("event", eventType.Name()),
		zap.Uint64("id", sid),
		zap.Time("start", time.Now()),
	)

	switch eventType {
	case types.EventTypeCoinBalanceChange:
		p.hd.AddSub(eventType.Name(), types.SubscriptionID(sid), p.hd.HandleBalanceChange)
	case types.EventTypePublish:
		p.hd.AddSub(eventType.Name(), types.SubscriptionID(sid), p.hd.HandlePublish)
	case types.EventTypeMoveEvent:
		p.hd.AddSub(eventType.Name(), types.SubscriptionID(sid), p.hd.HandleMove)
	case types.EventTypeNewObject:
		p.hd.AddSub(eventType.Name(), types.SubscriptionID(sid), p.hd.HandleNewObject)
	case types.EventTypeMutateObject:
		p.hd.AddSub(eventType.Name(), types.SubscriptionID(sid), p.hd.HandleMutateObject)
	case types.EventTypeDeleteObject:
		p.hd.AddSub(eventType.Name(), types.SubscriptionID(sid), p.hd.HandleDeleteObject)
	case types.EventTypeTransferObject:
		p.hd.AddSub(eventType.Name(), types.SubscriptionID(sid), p.hd.HandleTransferObject)
	default:
		err = fmt.Errorf("no handler for event: %s", eventType.Name())
	}
	return err
}
