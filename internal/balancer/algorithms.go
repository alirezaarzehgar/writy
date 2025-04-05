package balancer

import (
	"fmt"
	"math/rand"
	"sync/atomic"
)

type Algorithm[Client any] func(clients []Client) Client

type LoadBalancer[Client any] struct {
	algorithm Algorithm[Client]
}

func NewLoadBalancer[Client any](algorithm Algorithm[Client]) LoadBalancer[Client] {
	return LoadBalancer[Client]{algorithm: algorithm}
}

func (lb LoadBalancer[Client]) GetClient(rClients, wClients []Client) (Client, error) {
	if len(rClients) == 0 {
		if len(wClients) == 0 {
			var empty Client
			return empty, fmt.Errorf("no masters or slaves")
		}

		return lb.algorithm(wClients), nil
	}

	return lb.algorithm(rClients), nil
}

var rrCounter int64 = 1

func RoundRobin[Client any](clients []Client) Client {
	clen := int64(len(clients))
	client := clients[rrCounter-1]
	atomic.AddInt64(&rrCounter, 1)

	if rrCounter%clen == 0 {
		atomic.StoreInt64(&rrCounter, 1)
	}
	return client
}

func Randomized[Client any](clients []Client) Client {
	i := rand.Intn(len(clients))
	return clients[i]
}
