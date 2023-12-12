package app

import (
	"context"
	"net"
	"time"

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
	srv.grpcServer = grpc.NewServer(grpc.ChainUnaryInterceptor(srv.loggingInterceptor, srv.authInterceptor))
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

func (srv *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	var response pb.PingResponse
	return &response, srv.PingStorage(ctx)
}
