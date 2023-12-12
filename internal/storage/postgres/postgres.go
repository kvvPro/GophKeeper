package postgres

import (
	"context"
	"errors"
	_ "net/http/pprof"

	pb "github.com/kvvPro/gophkeeper/proto"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	ConnStr string
	pool    *pgxpool.Pool
}

func NewPSQLStr(ctx context.Context, connection string) (*PostgresStorage, error) {
	// init
	init := getInitQuery()
	pool, err := pgxpool.New(ctx, connection)
	if err != nil {
		return nil, err
	}

	//defer pool.Close()

	_, err = pool.Exec(ctx, init)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		ConnStr: connection,
		pool:    pool,
	}, nil
}

func (s *PostgresStorage) Quit(ctx context.Context) {
	s.pool.Close()
}

func (s *PostgresStorage) Ping(ctx context.Context) error {
	err := s.pool.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}

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

func (s *PostgresStorage) UpdateUserData(ctx context.Context, data *pb.UserData, owner string) (*pb.UserData, error) {
	var userInfo pb.UserData
	var userDataID string

	update := updateUserDataQuery()
	insertRes, err := s.pool.Exec(ctx, update, data.Login, data.Password, data.Key, owner)
	if err != nil {
		return nil, err
	}
	if insertRes.RowsAffected() == 0 {
		return nil, errors.New("user data not updated")
	}
	// update related metadata
	meta, err := s.UpdateMetadata(ctx, data.MetaInfo, userDataID)
	if err != nil {
		return nil, err
	}
	userInfo.MetaInfo = meta

	return &userInfo, nil
}

func (s *PostgresStorage) UpdateMetadata(ctx context.Context, data []*pb.Metadata, owner string) ([]*pb.Metadata, error) {
	var meta []*pb.Metadata

	// TODO
	return meta, nil
}

func (s *PostgresStorage) AddUserData(ctx context.Context, data *pb.UserData, owner string) error {

	return nil
}

func (s *PostgresStorage) GetMetadata(ctx context.Context, ownerID string) ([]*pb.Metadata, error) {

	metadata := []*pb.Metadata{}

	query := getMetadataQuery()
	result, err := s.pool.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	for result.Next() {
		var meta pb.Metadata
		err = result.Scan(&meta.Key, &meta.Value)
		if err != nil {
			return nil, err
		}
		metadata = append(metadata, &meta)
	}

	err = result.Err()
	if err != nil {
		return nil, err
	}

	return metadata, nil
}
