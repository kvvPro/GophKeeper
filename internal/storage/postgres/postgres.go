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

func (s *PostgresStorage) AddUser(ctx context.Context, user *pb.UserData) error {
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

func (s *PostgresStorage) GetUser(ctx context.Context, login string) (*pb.UserData, error) {
	var userInfo pb.UserData
	getUserQuery := getUserQuery()
	result := s.pool.QueryRow(ctx, getUserQuery, login)
	if err := result.Scan(&userInfo.Login, &userInfo.Password); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
