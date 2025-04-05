package balancer

import (
	"context"
	"fmt"
	"net"

	"github.com/alirezaarzehgar/writy/libwrity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type WrityService struct {
	readableConns []*grpc.ClientConn
	writableConns []*grpc.ClientConn
	libwrity.UnimplementedWrityServiceServer
}

func (ws *WrityService) Set(c context.Context, r *libwrity.SetRequest) (*libwrity.Empty, error) {
	return nil, nil
}

func (ws *WrityService) Get(c context.Context, r *libwrity.GetRequest) (*libwrity.GetResponse, error) {
	return nil, nil
}

func (ws *WrityService) Del(c context.Context, r *libwrity.DelRequest) (*libwrity.Empty, error) {
	return nil, nil
}

func (ws *WrityService) Keys(c context.Context, r *libwrity.KeysRequest) (*libwrity.KeysResponse, error) {
	return nil, nil
}

func (ws *WrityService) Flush(c context.Context, r *libwrity.Empty) (*libwrity.Empty, error) {
	return nil, nil
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
	RunningAddr       string
	ReflectionEnabled bool
	Slaves            StringArray
	Masters           StringArray
}

func Start(conf ServerConfig) error {
	writyService := WrityService{}

	for _, master := range conf.Masters {
		client, err := grpc.NewClient(master)
		if err != nil {
			return fmt.Errorf("failed to add new master connection: %s: %w", master, err)
		}
		writyService.readableConns = append(writyService.readableConns, client)
	}

	for _, slave := range conf.Slaves {
		client, err := grpc.NewClient(slave)
		if err != nil {
			return fmt.Errorf("failed to add new slave connection: %s: %w", slave, err)
		}
		writyService.readableConns = append(writyService.readableConns, client)
	}

	s := grpc.NewServer()

	libwrity.RegisterWrityServiceServer(s, &WrityService{})
	if conf.ReflectionEnabled {
		reflection.Register(s)
	}

	l, err := net.Listen("tcp", conf.RunningAddr)
	if err != nil {
		return fmt.Errorf("failed to create tcp connection: %w", err)
	}

	err = s.Serve(l)
	if err != nil {
		return fmt.Errorf("failed to serve gRPC connection: %w", err)
	}

	return nil
}
