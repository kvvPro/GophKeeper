package storage

import (
	"context"

	pb "github.com/kvvPro/gophkeeper/proto"
)

type Storage interface {
	Ping(ctx context.Context) error
	Quit(ctx context.Context)

	StartActionChain(ctx context.Context, clientID string) error
	CancelActionChain(ctx context.Context, clientID string) error
	DoneActionChain(ctx context.Context, clientID string) error
	ActionChainActive(ctx context.Context, clientID string) bool

	AddUser(ctx context.Context, user *pb.AuthInfo) error
	GetUser(ctx context.Context, login string) (*pb.AuthInfo, error)

	AddUserData(ctx context.Context, data *pb.UserData, owner string) error
	GetUserData(ctx context.Context, key string, owner string) (*pb.UserData, error)
	UpdateUserData(ctx context.Context, data *pb.UserData, owner string) error
	DeleteUserData(ctx context.Context, key string, owner string) error

	AddMetadata(ctx context.Context, data *pb.Metadata, userOwner string, keyData string) error
	UpdateMetadata(ctx context.Context, newData *pb.Metadata, userOwner string, keyData string) error
}
