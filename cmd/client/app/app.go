package app

import (
	"context"

	"github.com/kvvPro/gophkeeper/cmd/client/config"
	pb "github.com/kvvPro/gophkeeper/proto"
	"go.uber.org/zap"
)

var Sugar zap.SugaredLogger

type Client struct {
	// server address
	Address string
	// True if agent encrypts all messages
	needToEncrypt bool
	// Path to public key RSA
	publicKey string
	// auth token
	AuthToken string
}

// NewClient creates instance of client
func NewClient(settings *config.ClientFlags) (*Client, error) {
	return &Client{
		Address:       settings.Address,
		publicKey:     settings.CryptoKey,
		needToEncrypt: settings.CryptoKey != "",
	}, nil
}

func (cli *Client) AddUser(ctx context.Context, login string, password string) error {
	req := &pb.AuthInfo{
		Login:    login,
		Password: password,
	}
	token, err := cli.register(ctx, req)
	if err != nil {
		Sugar.Errorf("failed to add new user: %v", err)
		return err
	}

	cli.AuthToken = token
	return nil
}

func (cli *Client) Auth(ctx context.Context, login string, password string) error {
	req := &pb.AuthInfo{
		Login:    login,
		Password: password,
	}
	token, err := cli.auth(ctx, req)
	if err != nil {
		Sugar.Errorf("failed to pass authentication: %v", err)
		return err
	}

	cli.AuthToken = token
	return nil
}

func (cli *Client) WriteUserData(ctx context.Context, key string, login string, password string, meta map[string]string) error {
	metaInfo := make([]*pb.Metadata, 0)
	for k, v := range meta {
		metaInfo = append(metaInfo, &pb.Metadata{
			Key:   k,
			Value: v,
		})
	}
	req := &pb.UserData{
		Key:      key,
		Login:    login,
		Password: password,
		MetaInfo: metaInfo,
	}
	err := cli.putUserData(ctx, req)
	if err != nil {
		Sugar.Errorf("failed to put userdata: %v", err)
		return err
	}

	return nil
}

func (cli *Client) GetUserData(ctx context.Context, key string) (string, string, map[string]string, error) {
	data, err := cli.getUserData(ctx, key)
	if err != nil {
		Sugar.Errorf("failed to get userdata: %v", err)
		return "", "", nil, err
	}

	metaInfo := make(map[string]string, 0)
	for _, kv := range data.MetaInfo {
		metaInfo[kv.Key] = kv.Value
	}

	return data.Login, data.Password, metaInfo, nil
}
