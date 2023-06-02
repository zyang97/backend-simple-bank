package gapi

import (
	"fmt"

	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
	"github.com/techschool/simplebank/worker"
)

// Server serves gRPC requests for banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	pb.UnimplementedSimpleBankServer
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
