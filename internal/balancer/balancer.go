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

type WrityService struct {
	loadbalancer    LoadBalancer[libwrity.WrityServiceClient]
	readableClients []libwrity.WrityServiceClient
	writableClients []libwrity.WrityServiceClient
	libwrity.UnimplementedWrityServiceServer
}

func (ws *WrityService) Set(c context.Context, r *libwrity.SetRequest) (*libwrity.Empty, error) {
	for _, client := range ws.readableClients {
		go func(req *libwrity.SetRequest) {
			_, err := client.Set(context.TODO(), req)
			if err != nil {
				slog.Warn("failed to set value on slave", "error", err)
			}
		}(r)
		slog.Debug("balance to slave", "request", r)
	}
	for _, client := range ws.writableClients {
		_, err := client.Set(context.TODO(), r)
		if err != nil {
			slog.Warn("failed to set value on master", "request", r)
		}
		slog.Debug("balance to master", "request", r)
	}
	return &libwrity.Empty{}, nil
}

func (ws *WrityService) Get(c context.Context, r *libwrity.GetRequest) (*libwrity.GetResponse, error) {
	client := ws.loadbalancer.GetClient()
	return client.Get(c, r)
}

func (ws *WrityService) Del(c context.Context, r *libwrity.DelRequest) (*libwrity.Empty, error) {
	for _, client := range ws.readableClients {
		go client.Del(c, r)
	}
	for _, client := range ws.writableClients {
		client.Del(c, r)
	}
	return &libwrity.Empty{}, nil
}

func (ws *WrityService) Keys(c context.Context, r *libwrity.KeysRequest) (*libwrity.KeysResponse, error) {
	client := ws.loadbalancer.GetClient()
	return client.Keys(c, r)
}

func (ws *WrityService) Flush(c context.Context, r *libwrity.Empty) (*libwrity.Empty, error) {
	for _, client := range append(ws.readableClients, ws.writableClients...) {
		client.Flush(c, r)
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
	writyService := &WrityService{}

	for _, master := range conf.Masters {
		conn, err := grpc.NewClient(master, grpc.WithInsecure())
		if err != nil {
			return fmt.Errorf("failed to add new master connection: %s: %w", master, err)
		}
		client := libwrity.NewWrityServiceClient(conn)
		writyService.writableClients = append(writyService.writableClients, client)
	}

	for _, slave := range conf.Slaves {
		conn, err := grpc.NewClient(slave, grpc.WithInsecure())
		if err != nil {
			return fmt.Errorf("failed to add new slave connection: %s: %w", slave, err)
		}
		client := libwrity.NewWrityServiceClient(conn)
		writyService.readableClients = append(writyService.readableClients, client)
	}

	writyService.loadbalancer = NewLoadBalancer(writyService.readableClients, RoundRobin)

	s := grpc.NewServer()

	libwrity.RegisterWrityServiceServer(s, writyService)
	if conf.ReflectionEnabled {
		reflection.Register(s)
	}

	l, err := net.Listen("tcp", conf.RunningAddr)
	if err != nil {
		return fmt.Errorf("failed to create tcp connection: %w", err)
	}

	slog.Info("loadbalancer is ready")
	err = s.Serve(l)
	if err != nil {
		return fmt.Errorf("failed to serve gRPC connection: %w", err)
	}

	return nil
}
