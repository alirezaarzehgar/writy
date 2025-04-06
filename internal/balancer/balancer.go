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
	ReplicaAddresses       StringArray
	LoadBalancingAlgorithm Algorithm `json:"-"`
}

func Start(conf ServerConfig) error {
	balancerService := &LoadBalancerService{}

	for _, address := range conf.ReplicaAddresses {
		conn, err := grpc.NewClient(address, grpc.WithInsecure())
		if err != nil {
			return fmt.Errorf("failed to add new replica connection: %s: %w", address, err)
		}
		client := libwrity.NewWrityServiceClient(conn)
		replica := Replica{Connection: conn, Client: client, Address: address}
		balancerService.replicas = append(balancerService.replicas, replica)
	}

	balancerService.loadbalancer = NewLoadBalancer(conf.LoadBalancingAlgorithm)

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
