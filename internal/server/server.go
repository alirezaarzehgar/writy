package server

import (
	"context"

	"github.com/alirezaarzehgar/writy/libwrity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WrityService struct {
	libwrity.UnimplementedWrityServiceServer
}

func (WrityService) Set(context.Context, *libwrity.SetRequest) (*libwrity.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}

func (WrityService) Get(context.Context, *libwrity.GetRequest) (*libwrity.GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}

func (WrityService) Del(context.Context, *libwrity.DelRequest) (*libwrity.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Del not implemented")
}

func (WrityService) Keys(context.Context, *libwrity.KeysRequest) (*libwrity.KeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Keys not implemented")
}

func (WrityService) Flush(context.Context, *libwrity.Empty) (*libwrity.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Flush not implemented")
}

func New() {
	s := grpc.NewServer()
	libwrity.RegisterWrityServiceServer(s, &WrityService{})
}
