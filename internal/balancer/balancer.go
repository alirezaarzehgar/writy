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
	Replicas               StringArray
	LoadBalancingAlgorithm Algorithm[libwrity.WrityServiceClient] `json:"-"`
}

func Start(conf ServerConfig) error {
	balancerService := &LoadBalancerService{}

	for _, replica := range conf.Replicas {
		conn, err := grpc.NewClient(replica, grpc.WithInsecure())
		if err != nil {
			return fmt.Errorf("failed to add new replica connection: %s: %w", replica, err)
		}
		client := libwrity.NewWrityServiceClient(conn)
		balancerService.clients = append(balancerService.clients, client)
	}

	balancerService.loadbalancer = NewLoadBalancer[libwrity.WrityServiceClient](RoundRobin)

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
