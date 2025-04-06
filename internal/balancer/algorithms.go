package balancer

import (
	"fmt"
	"log/slog"
	"math/rand"
	"sync/atomic"
)

type Algorithm func(replicas []Replica) Replica

type LoadBalancer struct {
	algorithm Algorithm
}

func NewLoadBalancer(algorithm Algorithm) LoadBalancer {
	return LoadBalancer{algorithm: algorithm}
}

func (lb LoadBalancer) GetClient(replicas []Replica) (Replica, error) {
	var availableReplicas []Replica
	for _, r := range replicas {
		if !r.IsDown {
			availableReplicas = append(availableReplicas, r)
		}
	}

	if len(availableReplicas) == 0 {
		return Replica{}, fmt.Errorf("no masters or slaves")
	}

	r := lb.algorithm(availableReplicas)
	slog.Debug("balancing load", "address", r.Address)
	return r, nil
}

var rrCounter int64 = 0

func RoundRobin(replicas []Replica) Replica {
	clen := int64(len(replicas))

	if rrCounter >= clen {
		atomic.StoreInt64(&rrCounter, 0)
	}

	client := replicas[rrCounter]

	atomic.AddInt64(&rrCounter, 1)

	return client
}

func Randomized(replicas []Replica) Replica {
	i := rand.Intn(len(replicas))
	return replicas[i]
}
