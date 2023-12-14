package app

import (
	"context"

	pb "github.com/kvvPro/gophkeeper/proto"
	"google.golang.org/grpc"

	"github.com/kvvPro/gophkeeper/cmd/server/config"
	"github.com/kvvPro/gophkeeper/internal/storage"
	"github.com/kvvPro/gophkeeper/internal/storage/postgres"
	"go.uber.org/zap"
)

// Main logger
var Sugar zap.SugaredLogger

type Server struct {
	// Main storage for metrics, can be memstorage or postgresql type
	storage storage.Storage
	// Address of web server, where app woul be deployed
	// Format - Host[:Port]
	Address string
	// String connection to postgres DB
	// Format - "user=<user> password=<pass> host=<host> port=<port> dbname=<db> sslmode=<true/false>"
	DBConnection string
	// Key for decrypt and encrypt body of requests
	HashKey string
	// True if server would validate hash of all incoming requests
	CheckHash bool
	// True if server accepts encrypted messages from agent
	UseEncryption bool
	// Path to private key RSA
	PrivateKeyPath string
	// wait group for async saving
	//wg *sync.WaitGroup
	// implement GRPC server
	pb.UnimplementedExchangeServer
	// GRPC server
	grpcServer *grpc.Server
}

// NewServer creates app instance
func NewServer(settings *config.ServerFlags) (*Server, error) {

	newdb, err := postgres.NewPSQLStr(context.Background(), settings.DBConnection)
	if err != nil {
		return nil, err
	}

	return &Server{
		storage:        newdb,
		Address:        settings.Address,
		DBConnection:   settings.DBConnection,
		HashKey:        settings.HashKey,
		CheckHash:      settings.HashKey != "",
		PrivateKeyPath: settings.CryptoKey,
		UseEncryption:  settings.CryptoKey != "",
	}, nil
}

func (srv *Server) StartServer(ctx context.Context, srvFlags *config.ServerFlags) {

	// записываем в лог, что сервер запускается
	Sugar.Infow(
		"Starting server",
		"srvFlags", srvFlags,
	)

	srv.startGRPCServer(ctx)
}

func (srv *Server) StopServer(ctx context.Context) {
	srv.storage.Quit(ctx)
	srv.stopGRPCServer(ctx)
}

func (srv *Server) PingStorage(ctx context.Context) error {
	return srv.storage.Ping(ctx)
}
