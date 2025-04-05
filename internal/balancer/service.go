package balancer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/alirezaarzehgar/writy/libwrity"
	"google.golang.org/grpc"
)

type LoadBalancerService struct {
	loadbalancer    LoadBalancer[libwrity.WrityServiceClient]
	readableClients []libwrity.WrityServiceClient
	writableClients []libwrity.WrityServiceClient
	libwrity.UnimplementedLoadBalancerServiceServer
}

func (lbs *LoadBalancerService) Set(c context.Context, r *libwrity.SetRequest) (*libwrity.Empty, error) {
	for _, client := range lbs.readableClients {
		go func(req *libwrity.SetRequest) {
			_, err := client.Set(c, req)
			if err != nil {
				slog.Warn("failed to set value on slave", "error", err)
			}
		}(r)
	}
	for _, client := range lbs.writableClients {
		_, err := client.Set(c, r)
		if err != nil {
			slog.Warn("failed to set value on master", "request", r)
		}
	}
	return &libwrity.Empty{}, nil
}

func (lbs *LoadBalancerService) Get(c context.Context, r *libwrity.GetRequest) (*libwrity.GetResponse, error) {
	client, err := lbs.loadbalancer.GetClient(lbs.readableClients, lbs.writableClients)
	if err != nil {
		return nil, fmt.Errorf("there is no writy node")
	}
	return client.Get(c, r)
}

func (lbs *LoadBalancerService) Del(c context.Context, r *libwrity.DelRequest) (*libwrity.Empty, error) {
	for _, client := range lbs.readableClients {
		go func(req *libwrity.DelRequest) {
			_, err := client.Del(c, req)
			if err != nil {
				slog.Warn("failed to del key on slave", "error", err)
			}
		}(r)
	}
	for _, client := range lbs.writableClients {
		_, err := client.Del(c, r)
		if err != nil {
			slog.Warn("failed to delete value on master", "request", r)
		}
	}
	return &libwrity.Empty{}, nil
}

func (lbs *LoadBalancerService) Keys(c context.Context, r *libwrity.KeysRequest) (*libwrity.KeysResponse, error) {
	client, err := lbs.loadbalancer.GetClient(lbs.readableClients, lbs.writableClients)
	if err != nil {
		return nil, fmt.Errorf("there is no writy node")
	}
	return client.Keys(c, r)
}

func (lbs *LoadBalancerService) Flush(c context.Context, r *libwrity.Empty) (*libwrity.Empty, error) {
	for _, client := range append(lbs.readableClients, lbs.writableClients...) {
		_, err := client.Flush(c, r)
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

	if r.IsSlave {
		lbs.readableClients = append(lbs.readableClients, client)
	} else {
		lbs.writableClients = append(lbs.writableClients, client)
	}
	return &libwrity.Empty{}, nil
}
