package server

import (
	"context"

	"github.com/alirezaarzehgar/writy/libwrity"
	"google.golang.org/grpc"
)

type WrityService struct {
	libwrity.UnimplementedWrityServiceServer
}

func (w WrityService) Get(c context.Context, req *libwrity.GetRequest) (*libwrity.GetResponse, error) {
	return nil, nil
}

func (w WrityService) Set(c context.Context, req *libwrity.SetRequest) (*libwrity.SetResponse, error) {
	return nil, nil
}

func New() {
	s := grpc.NewServer()
	libwrity.RegisterWrityServiceServer(s, &WrityService{})
}
