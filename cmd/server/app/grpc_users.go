package app

import (
	"context"

	"github.com/kvvPro/gophkeeper/cmd/server/auth"
	pb "github.com/kvvPro/gophkeeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *Server) Register(ctx context.Context, req *pb.AuthInfo) (*pb.AuthResponse, error) {
	var response pb.AuthResponse
	// add user
	err := srv.AddUser(ctx, req)
	if err != nil {
		Sugar.Errorf("ошибка при добавлении нового пользователя: %v", err.Error())
		return nil, err
	}

	// create auth token
	uuid := auth.GenerateUUID()
	authToken, err := auth.BuildJWTString(req.Login, uuid)
	if err != nil {
		Sugar.Errorf("ошибка при генерации токена: %v", err.Error())
		return nil, err
	}

	response.AuthToken = authToken

	return &response, nil
}

func (srv *Server) Auth(ctx context.Context, req *pb.AuthInfo) (*pb.AuthResponse, error) {
	var response pb.AuthResponse
	// get user
	user, err := srv.GetUser(ctx, req)
	if err != nil {
		Sugar.Errorf("ошибка при получении пользователя: %v", err.Error())
		return nil, err
	}

	// check passwords
	if user.Password != req.Password {
		Sugar.Errorf("ошибка при аутентификации: %v", status.Errorf(codes.Unauthenticated, "wrong password"))
		return nil, err
	}

	// create auth token
	uuid := auth.GenerateUUID()
	authToken, err := auth.BuildJWTString(req.Login, uuid)
	if err != nil {
		Sugar.Errorf("ошибка при генерации токена: %v", err.Error())
		return nil, err
	}

	response.AuthToken = authToken

	return &response, nil
}
