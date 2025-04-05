package balancer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/alirezaarzehgar/writy/libwrity"
	"google.golang.org/grpc"
)

type LoadBalancerService struct {
	loadbalancer LoadBalancer[libwrity.WrityServiceClient]
	clients      []libwrity.WrityServiceClient
	libwrity.UnimplementedLoadBalancerServiceServer
}

func (lbs *LoadBalancerService) Set(c context.Context, r *libwrity.SetRequest) (*libwrity.Empty, error) {
	for _, client := range lbs.clients {
		go func(req *libwrity.SetRequest) {
			_, err := client.Set(context.TODO(), req)
			if err != nil {
				slog.Warn("failed to set value", "request", r, "error", err)
			}
		}(r)
	}
	return &libwrity.Empty{}, nil
}

func (lbs *LoadBalancerService) Get(c context.Context, r *libwrity.GetRequest) (*libwrity.GetResponse, error) {
	client, err := lbs.loadbalancer.GetClient(lbs.clients)
	if err != nil {
		return nil, fmt.Errorf("there is no writy node")
	}
	return client.Get(c, r)
}

func (lbs *LoadBalancerService) Del(c context.Context, r *libwrity.DelRequest) (*libwrity.Empty, error) {
	for _, client := range lbs.clients {
		go func(req *libwrity.DelRequest) {
			_, err := client.Del(context.TODO(), req)
			if err != nil {
				slog.Warn("failed to del key on replicas", "error", err)
			}
		}(r)
	}
	return &libwrity.Empty{}, nil
}

func (lbs *LoadBalancerService) Keys(c context.Context, r *libwrity.KeysRequest) (*libwrity.KeysResponse, error) {
	client, err := lbs.loadbalancer.GetClient(lbs.clients)
	if err != nil {
		return nil, fmt.Errorf("there is no writy node")
	}
	return client.Keys(context.TODO(), r)
}

func (lbs *LoadBalancerService) Flush(c context.Context, r *libwrity.Empty) (*libwrity.Empty, error) {
	for _, client := range lbs.clients {
		_, err := client.Flush(context.TODO(), r)
		if err != nil {
			slog.Warn("failed to flush node")
		}
	}
	return &libwrity.Empty{}, nil
}

func (lbs *LoadBalancerService) AddNode(c context.Context, r *libwrity.AddNodeRequest) (*libwrity.Empty, error) {
	conn, err := grpc.NewClient(r.Address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to add new master connection: %s: %w", r.Address, err)
	}
	client := libwrity.NewWrityServiceClient(conn)
	lbs.clients = append(lbs.clients, client)
	return &libwrity.Empty{}, nil
}
