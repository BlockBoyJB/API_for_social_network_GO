package authgrpc

import (
	"API_for_SN_go/internal/service"
	pb "API_for_SN_go/proto/auth"
	"context"
	"google.golang.org/grpc"
)

type authGrpc struct {
	pb.UnimplementedAuthServer
	authService service.Auth
}

func NewAuthGrpc(g *grpc.Server, authService service.Auth) {
	pb.RegisterAuthServer(g, &authGrpc{authService: authService})
}

func (g *authGrpc) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.SignInResponse, error) {
	token, err := g.authService.CreateToken(ctx, service.UserAuthInput{
		Username: in.Username,
		Password: in.Password,
	})
	if err != nil {
		return nil, err
	}

	return &pb.SignInResponse{Token: token}, nil
}

func (g *authGrpc) RefreshToken(ctx context.Context, in *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	token, err := g.authService.RefreshToken(ctx, in.Token)
	if err != nil {
		return nil, err
	}
	return &pb.RefreshTokenResponse{Token: token}, nil
}
