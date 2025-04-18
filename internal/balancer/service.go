package balancer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/alirezaarzehgar/writy/libwrity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	DefaultAssumeAliveCycle = time.Minute * 2
)

type Replica struct {
	Connection *grpc.ClientConn
	Client     libwrity.WrityServiceClient
	Address    string
	IsDown     bool
}

type LoadBalancerService struct {
	loadbalancer LoadBalancer
	replicas     []Replica
	libwrity.UnimplementedLoadBalancerServiceServer
}

func (lbs *LoadBalancerService) Set(c context.Context, req *libwrity.SetRequest) (*libwrity.Empty, error) {
	return writeOpFaileover(lbs, func(r Replica) error {
		_, err := r.Client.Set(context.TODO(), req)
		return err
	})
}

func (lbs *LoadBalancerService) Get(c context.Context, req *libwrity.GetRequest) (*libwrity.GetResponse, error) {
	return readOpFaileover[*libwrity.GetResponse](lbs, func(r Replica) (*libwrity.GetResponse, error) {
		return r.Client.Get(context.TODO(), req)
	})
}

func (lbs *LoadBalancerService) Del(c context.Context, req *libwrity.DelRequest) (*libwrity.Empty, error) {
	return writeOpFaileover(lbs, func(r Replica) error {
		_, err := r.Client.Del(context.TODO(), req)
		return err
	})
}

func (lbs *LoadBalancerService) Keys(c context.Context, req *libwrity.Empty) (*libwrity.KeysResponse, error) {
	return readOpFaileover[*libwrity.KeysResponse](lbs, func(r Replica) (*libwrity.KeysResponse, error) {
		return r.Client.Keys(context.TODO(), req)
	})
}

func (lbs *LoadBalancerService) Flush(c context.Context, req *libwrity.Empty) (*libwrity.Empty, error) {
	return writeOpFaileover(lbs, func(r Replica) error {
		_, err := r.Client.Flush(context.TODO(), req)
		return err
	})
}

func (lbs *LoadBalancerService) AddNode(c context.Context, r *libwrity.AddNodeRequest) (*libwrity.Empty, error) {
	conn, err := grpc.NewClient(r.Address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to add new master connection: %s: %w", r.Address, err)
	}
	client := libwrity.NewWrityServiceClient(conn)
	lbs.replicas = append(lbs.replicas, Replica{Connection: conn, Client: client, Address: r.Address})
	return &libwrity.Empty{}, nil
}

func (lbs *LoadBalancerService) DelNode(c context.Context, r *libwrity.DelNodeRequest) (*libwrity.Empty, error) {
	for i, r := range lbs.replicas {
		if r.Address == r.Address {
			lbs.replicas = append(lbs.replicas[:i], lbs.replicas[i+1:]...)
		}
	}
	return &libwrity.Empty{}, nil
}

func (lbs *LoadBalancerService) Nodes(c context.Context, r *libwrity.Empty) (*libwrity.NodesResponse, error) {
	var nodes []*libwrity.Node
	for _, r := range lbs.replicas {
		nodes = append(nodes, &libwrity.Node{Address: r.Address, Available: !r.IsDown})
	}
	return &libwrity.NodesResponse{Nodes: nodes}, nil
}

func readOpFaileover[T any](lbs *LoadBalancerService, task func(Replica) (T, error)) (T, error) {
	var empty T
	for {
		replica, err := lbs.loadbalancer.GetClient(lbs.replicas)
		if err != nil {
			return empty, fmt.Errorf("there is no writy node")
		}

		res, err := task(replica)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.Unavailable {
				slog.Debug("replica is unavailable", "error", err)
				markReplicaAsDown(lbs, replica)
				continue
			}
			return empty, err
		}
		return res, nil
	}
}

func writeOpFaileover(lbs *LoadBalancerService, task func(Replica) error) (*libwrity.Empty, error) {
	for {
		var availableReplicas []Replica
		for _, r := range lbs.replicas {
			if !r.IsDown {
				availableReplicas = append(availableReplicas, r)
			}
		}
		if len(availableReplicas) == 0 {
			return &libwrity.Empty{}, fmt.Errorf("there is no writy node")
		}

		master := availableReplicas[0]
		err := task(master)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.Unavailable {
				slog.Debug("replica is unavailable", "error", err)
				markReplicaAsDown(lbs, master)
				continue
			}
			return &libwrity.Empty{}, err
		}

		for _, r := range availableReplicas[1:] {
			go func(lbs *LoadBalancerService, replica Replica) {
				if err := task(replica); err != nil {
					slog.Debug("failed on write operation", "error", err)
					markReplicaAsDown(lbs, replica)
				}
			}(lbs, r)
		}

		return &libwrity.Empty{}, nil
	}
}

func markReplicaAsDown(lbs *LoadBalancerService, replica Replica) {
	for i, r := range lbs.replicas {
		if r == replica {
			lbs.replicas[i].IsDown = true
			go func(index int) {
				time.Sleep(DefaultAssumeAliveCycle)
				lbs.replicas[i].IsDown = false
			}(i)
			break
		}
	}
}
