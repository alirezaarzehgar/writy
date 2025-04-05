package balancer

import (
	"log/slog"
	"math/rand"
	"sync/atomic"
)

type Algorithm[Client any] func(clients []Client) Client

type LoadBalancer[Client any] struct {
	clients   []Client
	algorithm Algorithm[Client]
}

func NewLoadBalancer[Client any](clients []Client, algorithm Algorithm[Client]) LoadBalancer[Client] {
	return LoadBalancer[Client]{clients: clients, algorithm: algorithm}
}

func (lb LoadBalancer[Client]) GetClient() Client {
	slog.Debug("balance to a client")
	return lb.algorithm(lb.clients)
}

var rrCounter int64 = 0

func RoundRobin[Client any](clients []Client) Client {
	clen := int64(len(clients) - 1)
	client := clients[rrCounter]
	atomic.AddInt64(&rrCounter, 1)

	if rrCounter%clen == 0 {
		atomic.StoreInt64(&rrCounter, 0)
	}
	return client
}

func Randomized[Client any](clients []Client) Client {
	i := rand.Intn(len(clients))
	return clients[i]
}
