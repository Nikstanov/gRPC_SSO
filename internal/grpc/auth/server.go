package auth

import (
	"context"
	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	ssov1 "sso_service_grps/protos/gen/go/sso"
)

const (
	emptyValue = 0
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int32) (token string, err error)
	Register(ctx context.Context, email string, password string) (userID int64, err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func (s *serverAPI) Register(ctx context.Context, request *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	res := govalidator.IsEmail(request.GetEmail())
	if res != true {
		return nil, status.Error(codes.InvalidArgument, "Invalid email")
	}
	res = govalidator.StringLength(request.Password, "8", "15")
	if res != true {
		return nil, status.Error(codes.InvalidArgument, "Invalid password")
	}

	userID, err := s.auth.Register(ctx, request.GetEmail(), request.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) Login(ctx context.Context, request *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	res := govalidator.IsEmail(request.GetEmail())
	if res != true {
		return nil, status.Error(codes.InvalidArgument, "Invalid email")
	}
	if request.GetAppId() == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "Invalid AppId")
	}
	res = govalidator.StringLength(request.Password, "8", "15")
	if res != true {
		return nil, status.Error(codes.InvalidArgument, "Invalid password")
	}
	token, err := s.auth.Login(ctx, request.GetEmail(), request.GetPassword(), request.GetAppId())
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) mustEmbedUnimplementedAuthServer() {
	//TODO implement me
	panic("implement me")
}

func RegisterServerAPI(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})

}
