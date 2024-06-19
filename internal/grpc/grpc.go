package grpc

import (
	authgrpc "API_for_SN_go/internal/grpc/auth"
	"API_for_SN_go/internal/service"
	"google.golang.org/grpc"
)

func NewGRPC(services *service.Services) *grpc.Server {
	g := grpc.NewServer()
	authgrpc.NewAuthGrpc(g, services.Auth)
	return g
}
