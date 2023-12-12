package app

import (
	"context"

	pb "github.com/kvvPro/gophkeeper/proto"
)

func (srv *Server) PutUserData(ctx context.Context, req *pb.UserData) (*pb.PutDelInfoResponse, error) {
	var response pb.PutDelInfoResponse

	// TODO: add retry
	err := srv.ChangeUserData(ctx, req)
	if err != nil {
		return nil, err
	}

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
