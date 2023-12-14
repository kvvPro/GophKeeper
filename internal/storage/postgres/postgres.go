package postgres

import (
	"context"
	_ "net/http/pprof"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	ConnStr      string
	pool         *pgxpool.Pool
	transactions map[string]pgx.Tx
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

func (s *PostgresStorage) StartActionChain(ctx context.Context, clientID string) error {
	transaction, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	s.transactions[clientID] = transaction

	return nil
}

func (s *PostgresStorage) CancelActionChain(ctx context.Context, clientID string) error {
	err := s.transactions[clientID].Rollback(ctx)
	delete(s.transactions, clientID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) DoneActionChain(ctx context.Context, clientID string) error {
	return s.transactions[clientID].Commit(ctx)
}

func (s *PostgresStorage) ActionChainActive(ctx context.Context, clientID string) bool {
	_, ok := s.transactions[clientID]

	return ok
}
