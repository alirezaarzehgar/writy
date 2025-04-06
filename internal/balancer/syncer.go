package balancer

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"github.com/alirezaarzehgar/writy/libwrity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	DefaultSyncReplicasCycle = time.Second * 10
)

type Syncer struct {
	lbs *LoadBalancerService
}

func NewSyncer(lbs *LoadBalancerService) Syncer {
	return Syncer{lbs: lbs}
}

func (s Syncer) Start() {
	go func() {
		for {
			select {
			case <-time.NewTicker(DefaultSyncReplicasCycle).C:
				s.syncLoop()
			}
		}
	}()
}

func (s Syncer) syncLoop() {
	for {
		var availableReplicas []Replica
		for _, r := range s.lbs.replicas {
			if !r.IsDown {
				availableReplicas = append(availableReplicas, r)
			}
		}
		if len(availableReplicas) == 0 {
			slog.Debug("there is no replica to sync")
			return
		}

		master := availableReplicas[0]
		keys, err := master.Client.Keys(context.TODO(), &libwrity.Empty{})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.Unavailable {
				slog.Debug("replica is unavailable", "error", err)
				markReplicaAsDown(s.lbs, master)
				continue
			}
			slog.Debug("failed to get master keys", "error", err, "master address", master.Address)
			return
		}

		storage := make(map[string]string)
		for _, k := range keys.Keys {
			r, err := master.Client.Get(context.TODO(), &libwrity.GetRequest{Key: k})
			if err != nil {
				slog.Warn("failed to get data", "key", k, "error", err)
				continue
			}
			storage[k] = r.Value
		}

		for _, r := range availableReplicas[1:] {
			go s.sync(storage, keys.Keys, r)
		}
		return
	}
}

func (s Syncer) sync(masterStorage map[string]string, masterKeys []string, replica Replica) {
	keys, err := replica.Client.Keys(context.TODO(), &libwrity.Empty{})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unavailable {
			slog.Debug("replica is unavailable", "error", err)
			markReplicaAsDown(s.lbs, replica)
			return
		}

		for k, v := range masterStorage {
			_, err := replica.Client.Set(context.TODO(), &libwrity.SetRequest{Key: k, Value: v})
			if err != nil {
				slog.Warn("failed to set", "error", err, "key", k, "value", v)
			}
		}
		return
	}
	replicaKeys := keys.Keys

	for k, v := range masterStorage {
		if !slices.Contains(replicaKeys, k) {
			_, err := replica.Client.Set(context.TODO(), &libwrity.SetRequest{Key: k, Value: v})
			if err != nil {
				slog.Warn("failed to set", "error", err, "key", k, "value", v)
			}
		}
	}

	for _, k := range replicaKeys {
		if !slices.Contains(masterKeys, k) {
			_, err := replica.Client.Del(context.TODO(), &libwrity.DelRequest{Key: k})
			if err != nil {
				slog.Warn("failed to delete key", "error", err, "key", k)
			}
		}
	}
}
