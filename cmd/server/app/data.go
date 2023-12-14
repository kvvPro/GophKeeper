package app

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kvvPro/gophkeeper/cmd/server/auth"
	pb "github.com/kvvPro/gophkeeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *Server) ChangeUserData(ctx context.Context, userInfo *pb.UserData) error {
	// get login from context
	clientInfo := ctx.Value(ctxKey("clientInfo")).(*auth.ClientInfo)
	if clientInfo == nil {
		return status.Error(codes.Aborted, "not found current client info")
	}

	srv.storage.StartActionChain(ctx, clientInfo.ClientID)
	defer srv.storage.CancelActionChain(ctx, clientInfo.ClientID)
	// check if userdata alredy exists
	oldUserData, err := srv.storage.GetUserData(ctx, userInfo.Key, clientInfo.UserLogin)

	switch err {
	case pgx.ErrNoRows:
		// данных нет, добавляем новые
		err = srv.storage.AddUserData(ctx, userInfo, clientInfo.UserLogin)
	case nil:
		// данные уже есть, надо обновить
		err = srv.storage.UpdateUserData(ctx, userInfo, clientInfo.UserLogin)
	default:
		Sugar.Errorln(err)
		return err
	}

	if err != nil {
		Sugar.Errorln(err)
		return err
	}

	// update metadata
	// check for each pair what to do - update or add
	oldMetaMap := make(map[string]string, 0)
	for _, el := range oldUserData.MetaInfo {
		oldMetaMap[el.Key] = el.Value
	}

	for _, kv := range userInfo.MetaInfo {
		if _, ok := oldMetaMap[kv.Key]; !ok {
			// add new metadata
			err = srv.storage.AddMetadata(ctx, kv, clientInfo.UserLogin, userInfo.Key)
		} else {
			// update existed metadata
			err = srv.storage.UpdateMetadata(ctx, kv, clientInfo.UserLogin, userInfo.Key)
		}
	}

	if err != nil {
		return err
	}

	err = srv.storage.DoneActionChain(ctx, clientInfo.ClientID)
	return err
}

func (srv *Server) GetUserData(ctx context.Context, key string) (*pb.UserData, error) {
	// get login from context
	clientInfo := ctx.Value(ctxKey("clientInfo")).(*auth.ClientInfo)
	if clientInfo == nil {
		return nil, status.Error(codes.Aborted, "not found current client info")
	}

	data, err := srv.storage.GetUserData(ctx, key, clientInfo.UserLogin)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (srv *Server) DeleteUserData(ctx context.Context, key string) error {
	// get login from context
	clientInfo := ctx.Value(ctxKey("clientInfo")).(*auth.ClientInfo)
	if clientInfo == nil {
		return status.Error(codes.Aborted, "not found current client info")
	}

	err := srv.storage.DeleteUserData(ctx, key, clientInfo.UserLogin)
	if err != nil {
		return err
	}

	return nil
}
