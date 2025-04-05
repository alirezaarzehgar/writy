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

func (lb LoadBalancer[Client]) GetClient(clients []Client) (Client, error) {
	if len(clients) == 0 {
		var empty Client
		return empty, fmt.Errorf("no masters or slaves")
	}

	return lb.algorithm(clients), nil
}

var rrCounter int64 = 0

func RoundRobin[Client any](clients []Client) Client {
	clen := int64(len(clients))
	client := clients[rrCounter]
	atomic.AddInt64(&rrCounter, 1)

	if rrCounter == clen {
		atomic.StoreInt64(&rrCounter, 0)
	}
	return client
}

func Randomized[Client any](clients []Client) Client {
	i := rand.Intn(len(clients))
	return clients[i]
}
