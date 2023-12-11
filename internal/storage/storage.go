package storage

import (
	"context"

	pb "github.com/kvvPro/gophkeeper/proto"
)

type Storage interface {
	Ping(ctx context.Context) error
	Quit(ctx context.Context)
	AddUser(ctx context.Context, user *pb.UserData) error
	GetUser(ctx context.Context, login string) (*pb.UserData, error)
}
