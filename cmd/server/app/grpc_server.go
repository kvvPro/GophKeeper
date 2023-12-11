package app

import (
	"context"
	"net"
	"time"

	"github.com/kvvPro/gophkeeper/cmd/server/auth"
	pb "github.com/kvvPro/gophkeeper/proto"
	"google.golang.org/grpc"
)

func (srv *Server) startGRPCServer(ctx context.Context) {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", srv.Address)
	if err != nil {
		Sugar.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	srv.grpcServer = grpc.NewServer(grpc.ChainUnaryInterceptor(srv.loggingInterceptor))
	// регистрируем сервис
	pb.RegisterExchangeServer(srv.grpcServer, srv)

	Sugar.Infoln("Сервер gRPC начал работу")
	// получаем запрос gRPC
	go func() {
		if err := srv.grpcServer.Serve(listen); err != nil {
			Sugar.Fatalw(err.Error(), "event", "start grpc server")
		}
	}()
}

func (srv *Server) stopGRPCServer(ctx context.Context) {
	stopped := make(chan struct{})
	Sugar.Infoln("Попытка мягко завершить сервер")
	timeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	go func() {
		defer cancel()
		srv.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-timeout.Done():
		Sugar.Errorf("Ошибка при попытке мягко завершить grpc-сервер: %v", "timeout is expired")
		srv.grpcServer.Stop()
	}
}

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

func (srv *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	var response pb.PingResponse
	return &response, srv.PingStorage(ctx)
}

func (srv *Server) Register(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	var response pb.AuthResponse
	// add user
	err := srv.AddUser(ctx, req.UserInfo)
	if err != nil {
		Sugar.Errorf("ошибка при добавлении нового пользователя: %v", err.Error())
		return nil, err
	}

	// create auth token
	authToken, err := auth.BuildJWTString(req.UserInfo.Login)
	if err != nil {
		Sugar.Errorf("ошибка при генерации токена: %v", err.Error())
		return nil, err
	}

	response.AuthToken = authToken

	return &response, nil
}

func (srv *Server) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	var response pb.AuthResponse

	return &response, nil
}

func (srv *Server) PutUserData(ctx context.Context, req *pb.UserData) (*pb.PutDelInfoResponse, error) {
	var response pb.PutDelInfoResponse

	return &response, nil
}

func (srv *Server) PutRawData(ctx context.Context, req *pb.RawData) (*pb.PutDelInfoResponse, error) {
	var response pb.PutDelInfoResponse

	return &response, nil
}

func (srv *Server) PutTextData(ctx context.Context, req *pb.TextData) (*pb.PutDelInfoResponse, error) {
	var response pb.PutDelInfoResponse

	return &response, nil
}

func (srv *Server) PutCardData(ctx context.Context, req *pb.CardData) (*pb.PutDelInfoResponse, error) {
	var response pb.PutDelInfoResponse

	return &response, nil
}

func (srv *Server) DeleteInfo(ctx context.Context, req *pb.InfoRequest) (*pb.PutDelInfoResponse, error) {
	var response pb.PutDelInfoResponse

	return &response, nil
}

func (srv *Server) GetInfo(ctx context.Context, req *pb.InfoRequest) (*pb.GetInfoResponse, error) {
	var response pb.GetInfoResponse

	return &response, nil
}
