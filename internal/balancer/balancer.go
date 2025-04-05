package balancer

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/alirezaarzehgar/writy/libwrity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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
