package main

import (
	"context"
	"fmt"
	"net"

	"github.com/strahe/suialert/cmd"
)

type A struct {
	DialContext func(ctx context.Context, network, addr string) (net.Conn, error)
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}
