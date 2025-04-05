package balancer

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/alirezaarzehgar/writy/libwrity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

type StringArray []string

func (i *StringArray) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *StringArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type ServerConfig struct {
	RunningAddr            string
	ReflectionEnabled      bool
	Slaves                 StringArray
	Masters                StringArray
	LoadBalancingAlgorithm Algorithm[libwrity.WrityServiceClient] `json:"-"`
}

func Start(conf ServerConfig) error {
	balancerService := &LoadBalancerService{}

	for _, master := range conf.Masters {
		conn, err := grpc.NewClient(master, grpc.WithInsecure())
		if err != nil {
			return fmt.Errorf("failed to add new master connection: %s: %w", master, err)
		}
		client := libwrity.NewWrityServiceClient(conn)
		balancerService.writableClients = append(balancerService.writableClients, client)
	}

	for _, slave := range conf.Slaves {
		conn, err := grpc.NewClient(slave, grpc.WithInsecure())
		if err != nil {
			return fmt.Errorf("failed to add new slave connection: %s: %w", slave, err)
		}
		client := libwrity.NewWrityServiceClient(conn)
		balancerService.readableClients = append(balancerService.readableClients, client)
	}

	balancerService.loadbalancer = NewLoadBalancer(balancerService.readableClients, RoundRobin)

	s := grpc.NewServer()

	libwrity.RegisterLoadBalancerServiceServer(s, balancerService)
	if conf.ReflectionEnabled {
		reflection.Register(s)
	}

	l, err := net.Listen("tcp", conf.RunningAddr)
	if err != nil {
		return fmt.Errorf("failed to create tcp connection: %w", err)
	}

	slog.Info("loadbalancer is ready to work")
	err = s.Serve(l)
	if err != nil {
		return fmt.Errorf("failed to serve gRPC connection: %w", err)
	}

	return nil
}
