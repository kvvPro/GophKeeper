package postgres

import (
	"context"
	"errors"
	_ "net/http/pprof"

	pb "github.com/kvvPro/gophkeeper/proto"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *PostgresStorage) AddUser(ctx context.Context, user *pb.AuthInfo) error {
	addUserQuery := addUserQuery()
	insertRes, err := s.pool.Exec(ctx, addUserQuery, user.Login, user.Password)
	if err != nil {
		return err
	}
	if insertRes.RowsAffected() == 0 {
		return errors.New("can't add user")
	}

	return nil
}

func (s *PostgresStorage) GetUser(ctx context.Context, login string) (*pb.AuthInfo, error) {
	var userInfo pb.AuthInfo
	getUserQuery := getUserQuery()
	result := s.pool.QueryRow(ctx, getUserQuery, login)
	if err := result.Scan(&userInfo.Login, &userInfo.Password); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
