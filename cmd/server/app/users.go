package app

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kvvPro/gophkeeper/internal/retry"
	pb "github.com/kvvPro/gophkeeper/proto"
)

func (srv *Server) AddUser(ctx context.Context, userInfo *pb.AuthInfo) error {
	// add user
	err := retry.Do(func() error {
		return srv.storage.AddUser(ctx, userInfo)
	},
		retry.RetryIf(func(errAttempt error) bool {
			var pgErr *pgconn.PgError
			if errors.As(errAttempt, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				return true
			}
			return false
		}),
		retry.Attempts(3),
		retry.InitDelay(5*time.Millisecond),
		retry.Step(2*time.Millisecond),
		retry.Context(ctx),
	)

	if err != nil {
		Sugar.Errorln(err)
		return err
	}

	return err
}

func (srv *Server) GetUser(ctx context.Context, userInfo *pb.AuthInfo) (*pb.AuthInfo, error) {
	// get user
	var user *pb.AuthInfo
	var err error

	err = retry.Do(func() error {
		user, err = srv.storage.GetUser(ctx, userInfo.Login)
		return err
	},
		retry.RetryIf(func(errAttempt error) bool {
			var pgErr *pgconn.PgError
			if errors.As(errAttempt, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				return true
			}
			return false
		}),
		retry.Attempts(3),
		retry.InitDelay(5*time.Millisecond),
		retry.Step(2*time.Millisecond),
		retry.Context(ctx),
	)

	if err != nil {
		Sugar.Errorln(err)
		return nil, err
	}

	return user, err
}
