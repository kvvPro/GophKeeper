package app

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/kvvPro/gophkeeper/internal/model"
	pb "github.com/kvvPro/gophkeeper/proto"
)

func (cli *Client) connect(ctx context.Context) (*grpc.ClientConn, pb.ExchangeClient, error) {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(cli.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		Sugar.Fatal(err)
		return nil, nil, err
	}

	// получаем переменную интерфейсного типа ExchangeClient,
	// через которую будем отправлять сообщения
	c := pb.NewExchangeClient(conn)

	return conn, c, nil
}

func (cli *Client) register(ctx context.Context, req *pb.AuthInfo) (string, error) {
	conn, c, err := cli.connect(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	response, err := c.Register(ctx, req)
	if response == nil || err != nil {
		// smth is wrong
		return "", err
	}

	return response.AuthToken, nil
}

func (cli *Client) auth(ctx context.Context, req *pb.AuthInfo) (string, error) {
	conn, c, err := cli.connect(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	response, err := c.Auth(ctx, req)
	if response == nil || err != nil {
		// smth is wrong
		return "", err
	}

	return response.AuthToken, nil
}

func (cli *Client) putUserData(ctx context.Context, req *pb.UserData) error {
	conn, c, err := cli.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	md := metadata.New(map[string]string{"Authorization": cli.AuthToken})
	ctxClient := metadata.NewOutgoingContext(ctx, md)

	response, err := c.PutUserData(ctxClient, req)
	if response == nil || err != nil {
		// smth is wrong
		return err
	}

	return nil
}

func (cli *Client) getUserData(ctx context.Context, key string) (*pb.UserData, error) {
	conn, c, err := cli.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	md := metadata.New(map[string]string{"Authorization": cli.AuthToken})
	ctxClient := metadata.NewOutgoingContext(ctx, md)

	req := &pb.InfoRequest{
		Key:          key,
		ResourceType: model.DataTypeUserData,
	}
	response, err := c.GetInfo(ctxClient, req)
	if response == nil || err != nil {
		// smth is wrong
		return nil, err
	}

	return response.User, nil
}
