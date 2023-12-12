package storage

import (
	"context"

	pb "github.com/kvvPro/gophkeeper/proto"
)

type Storage interface {
	Ping(ctx context.Context) error
	Quit(ctx context.Context)
	AddUser(ctx context.Context, user *pb.AuthInfo) error
	GetUser(ctx context.Context, login string) (*pb.AuthInfo, error)
	GetUserData(ctx context.Context, key string, owner string) (*pb.UserData, error)
	UpdateUserData(ctx context.Context, data *pb.UserData, owner string) (*pb.UserData, error)
	AddUserData(ctx context.Context, data *pb.UserData, owner string) error
}
