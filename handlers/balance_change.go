package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/strahe/suialert/types"
	"go.uber.org/zap"
)

// HandleBalanceChange handle the balance change events.
func (e *SubHandler) HandleBalanceChange(ctx context.Context, event *types.Subscription) error {
	if event == nil {
		return nil
	}
	var er types.EventResult
	if err := json.Unmarshal(event.Result, &er); err != nil {
		return fmt.Errorf("error unmarshalling event result: %s", err.Error())
	}

	for name, raw := range er.Event {
		switch name {
		case types.EventCoinBalanceChange.Name():
			ed := types.CoinBalanceChange{}
			if err := json.Unmarshal(raw, &ed); err != nil {
				zap.S().Errorf("error unmarshalling coin balance change event: %s", err.Error())
				return err
			}
			zap.S().Info(ed)
		default:
			zap.S().Warnf("unknown event name: %s", name)
		}
	}
	return nil
}
