package postgres

import (
	"context"
	"errors"
	_ "net/http/pprof"

	pb "github.com/kvvPro/gophkeeper/proto"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *PostgresStorage) UpdateMetadata(ctx context.Context, newData *pb.Metadata, userOwner string, keyData string) error {
	update := updateMetadataQuery()
	updateRes, err := s.pool.Exec(ctx, update, newData.Key, newData.Value, userOwner, keyData)
	if err != nil {
		return err
	}
	if updateRes.RowsAffected() == 0 {
		return errors.New("can't update metadata")
	}

	return nil
}

func (s *PostgresStorage) AddMetadata(ctx context.Context, data *pb.Metadata, userOwner string, keyData string) error {
	addUserQuery := addMetadataQuery()
	insertRes, err := s.pool.Exec(ctx, addUserQuery, data.Key, data.Value, userOwner, keyData)
	if err != nil {
		return err
	}
	if insertRes.RowsAffected() == 0 {
		return errors.New("can't add user")
	}

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
