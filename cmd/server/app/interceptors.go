package app

import (
	"context"
	"time"

	"github.com/kvvPro/gophkeeper/cmd/server/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (srv *Server) loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	h, err := handler(ctx, req)

	duration := time.Since(start)

	Sugar.Infoln(
		"uri", info.FullMethod,
		"duration", duration,
		"err", err,
	)
	return h, err
}

type ctxKey string

func (srv *Server) authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var authToken string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		param := md.Get("Authorization")
		if len(param) > 0 {
			authToken = param[0]
		}
	}
	if len(authToken) == 0 {
		return nil, status.Error(codes.Aborted, "not found authorization token")
	}

	userInfo, err := auth.GetUserInfo(authToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	newContext := context.WithValue(ctx, ctxKey("login"), userInfo.Login)

	return handler(newContext, req)
}
