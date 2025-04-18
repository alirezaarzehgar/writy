package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/alirezaarzehgar/writy/internal/writy"
	"github.com/alirezaarzehgar/writy/libwrity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type WrityService struct {
	writy *writy.Writy
	libwrity.UnimplementedWrityServiceServer
}

func (ws *WrityService) Set(c context.Context, r *libwrity.SetRequest) (*libwrity.Empty, error) {
	if err := throwNullStrinError(r.Key, r.Value); err != nil {
		return nil, err
	}

	err := ws.writy.Set(r.Key, r.Value)
	slog.Info("set", "key", r.Key, "val", r.Value)
	return &libwrity.Empty{}, err
}

func (ws *WrityService) Get(c context.Context, r *libwrity.GetRequest) (*libwrity.GetResponse, error) {
	if err := throwNullStrinError(r.Key); err != nil {
		return nil, err
	}

	v := ws.writy.Get(r.Key)
	if v == nil {
		return nil, fmt.Errorf("not found: %s", r.Key)
	}
	return &libwrity.GetResponse{Value: fmt.Sprint(v)}, nil
}

func (ws *WrityService) Del(c context.Context, r *libwrity.DelRequest) (*libwrity.Empty, error) {
	if err := throwNullStrinError(r.Key); err != nil {
		return nil, err
	}

	return &libwrity.Empty{}, ws.writy.Del(r.Key)
}

func (ws *WrityService) Keys(c context.Context, r *libwrity.Empty) (*libwrity.KeysResponse, error) {
	keys := ws.writy.Keys()
	if keys == nil {
		return nil, fmt.Errorf("storage is empty now")
	}
	return &libwrity.KeysResponse{Keys: keys}, nil
}

func (ws *WrityService) Flush(c context.Context, r *libwrity.Empty) (*libwrity.Empty, error) {
	ws.writy.Cleanup()
	return &libwrity.Empty{}, nil
}

type ServerConfig struct {
	DbPath, RunningAddr string
	ReflectionEnabled   bool
	GcCycle             time.Duration
}

func Start(conf ServerConfig) error {
	s := grpc.NewServer()

	writy.DefaultGarbageCollectorCycle = conf.GcCycle

	w, err := writy.New(conf.DbPath)
	if err != nil {
		return fmt.Errorf("failed to create writy instalce: %w", err)
	}

	libwrity.RegisterWrityServiceServer(s, &WrityService{writy: w})
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

func throwNullStrinError(values ...string) error {
	for _, v := range values {
		if v == "" {
			return fmt.Errorf("empty string is not permitted")
		}
	}
	return nil
}
