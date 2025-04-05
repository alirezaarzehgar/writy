package balancer

import (
	"context"
	"log/slog"

	"github.com/alirezaarzehgar/writy/libwrity"
)

type LoadBalancerService struct {
	loadbalancer    LoadBalancer[libwrity.WrityServiceClient]
	readableClients []libwrity.WrityServiceClient
	writableClients []libwrity.WrityServiceClient
	libwrity.UnimplementedLoadBalancerServiceServer
}

func (ws *LoadBalancerService) Set(c context.Context, r *libwrity.SetRequest) (*libwrity.Empty, error) {
	for _, client := range ws.readableClients {
		go func(req *libwrity.SetRequest) {
			_, err := client.Set(c, req)
			if err != nil {
				slog.Warn("failed to set value on slave", "error", err)
			}
		}(r)
	}
	for _, client := range ws.writableClients {
		_, err := client.Set(c, r)
		if err != nil {
			slog.Warn("failed to set value on master", "request", r)
		}
	}
	return &libwrity.Empty{}, nil
}

func (ws *LoadBalancerService) Get(c context.Context, r *libwrity.GetRequest) (*libwrity.GetResponse, error) {
	client := ws.loadbalancer.GetClient()
	return client.Get(c, r)
}

func (ws *LoadBalancerService) Del(c context.Context, r *libwrity.DelRequest) (*libwrity.Empty, error) {
	for _, client := range ws.readableClients {
		go func(req *libwrity.DelRequest) {
			_, err := client.Del(c, req)
			if err != nil {
				slog.Warn("failed to del key on slave", "error", err)
			}
		}(r)
	}
	for _, client := range ws.writableClients {
		_, err := client.Del(c, r)
		if err != nil {
			slog.Warn("failed to delete value on master", "request", r)
		}
	}
	return &libwrity.Empty{}, nil
}

func (ws *LoadBalancerService) Keys(c context.Context, r *libwrity.KeysRequest) (*libwrity.KeysResponse, error) {
	client := ws.loadbalancer.GetClient()
	return client.Keys(c, r)
}

func (ws *LoadBalancerService) Flush(c context.Context, r *libwrity.Empty) (*libwrity.Empty, error) {
	for _, client := range append(ws.readableClients, ws.writableClients...) {
		_, err := client.Flush(c, r)
		if err != nil {
			slog.Warn("failed to flush node")
		}
	}
	return &libwrity.Empty{}, nil
}
