package app

import (
	"context"

	"github.com/jackc/pgx/v5"
	pb "github.com/kvvPro/gophkeeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *Server) ChangeUserData(ctx context.Context, userInfo *pb.UserData) error {
	// get login from context
	login := ctx.Value(ctxKey("login")).(string)
	if len(login) == 0 {
		return status.Error(codes.Aborted, "not found current user login")
	}

	// check if userdata alredy exists
	_, err := srv.storage.GetUserData(ctx, userInfo.Key, login)

	switch err {
	case pgx.ErrNoRows:
		// данных нет, добавляем новые
		err = srv.storage.AddUserData(ctx, userInfo, login)
	case nil:
		// данные уже есть, надо обновить
		_, err = srv.storage.UpdateUserData(ctx, userInfo, login)
	default:
		Sugar.Errorln(err)
		return err
	}

	if err != nil {
		Sugar.Errorln(err)
		return err
	}

	return err
}
