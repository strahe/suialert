package client

import (
	"context"

	"github.com/strahe/suialert/types"

	"github.com/filecoin-project/go-jsonrpc"
)

// Client is the client interface for interacting with the sui node.
type Client struct {
	GetEvents        func(ctx context.Context, query types.EventQuery, cursor *types.EventID, limit uint, descendingOrder bool) (*types.EventPage, error)
	SubscribeEvent   func(ctx context.Context, query types.SubscribeEventQuery) (uint64, error)
	UnsubscribeEvent func(ctx context.Context, id uint64) (bool, error)
}

type SubscribeEventHandler interface {
	SubscribeEvent(ctx context.Context, r jsonrpc.RawParams) error
}

// NewClient creates a new client.
func NewClient(ctx context.Context, addr string, hd SubscribeEventHandler) (*Client, func(), error) {

	rpcOpts := []jsonrpc.Option{
		jsonrpc.WithClientHandler("Sui", hd),
		jsonrpc.WithClientHandlerAlias("sui_subscribeEvent", "Sui.SubscribeEvent"),
	}

	var client Client
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Sui", []interface{}{&client}, nil, rpcOpts...)
	if err != nil {
		return nil, nil, err
	}

	return &client, closer, nil
}
