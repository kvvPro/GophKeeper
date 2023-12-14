package postgres

import (
	"context"
	"errors"
	_ "net/http/pprof"

	"github.com/kvvPro/gophkeeper/internal/model"
	pb "github.com/kvvPro/gophkeeper/proto"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *PostgresStorage) GetUserData(ctx context.Context, key string, owner string) (*pb.UserData, error) {
	var userInfo pb.UserData
	var userDataID string
	getUserQuery := getUserDataQuery()
	result := s.pool.QueryRow(ctx, getUserQuery, key, owner)
	if err := result.Scan(&userDataID, &userInfo.Login, &userInfo.Password); err != nil {
		return nil, err
	}
	metadata, err := s.GetMetadata(ctx, userDataID)
	if err != nil {
		return nil, err
	}
	userInfo.MetaInfo = metadata

	return &userInfo, nil
}

func (s *PostgresStorage) UpdateUserData(ctx context.Context, data *pb.UserData, owner string) error {
	update := updateUserDataQuery()
	updateRes, err := s.pool.Exec(ctx, update, data.Login, data.Password, data.Key, owner)
	if err != nil {
		return err
	}
	if updateRes.RowsAffected() == 0 {
		return errors.New("can't update user data")
	}

	return nil
}

func (s *PostgresStorage) AddUserData(ctx context.Context, data *pb.UserData, owner string) error {
	addUserDataQuery := addUserDataQuery()
	insertRes, err := s.pool.Exec(ctx, addUserDataQuery, owner, data.Key, model.DataTypeUserData, data.Login, data.Password)
	if err != nil {
		return err
	}
	if insertRes.RowsAffected() == 0 {
		return errors.New("can't add user data")
	}

	return nil
}

func (s *PostgresStorage) DeleteUserData(ctx context.Context, key string, owner string) error {
	deleteUserDataQuery := deleteDataQuery()
	deleteRes, err := s.pool.Exec(ctx, deleteUserDataQuery, key, owner)
	if err != nil {
		return err
	}
	if deleteRes.RowsAffected() == 0 {
		return errors.New("can't delete user data")
	}

	return nil
}
