package tinyrpcgo

import (
	"net/rpc"
)

type Client struct {
	*rpc.Client
}

// Functional Options Pattern
