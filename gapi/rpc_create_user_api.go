package gapi

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/util"
	"github.com/techschool/simplebank/val"
	"github.com/techschool/simplebank/worker"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Validate create user request.
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	// Create user follows request.
	hashedPassword, err := util.HashedPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			// Send verify email.
			taskPayload := worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}

			options := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(time.Second * 10),
				asynq.Queue(worker.QUEUECRITICAL),
			}

			return server.taskDistributor.DistributorTaskSendVerifyEmail(ctx, &taskPayload, options...)
		},
	}

	txResult, err := server.store.CreateUserTx(ctx, arg)

	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Unimplemented, "failed to create user: %s", err)
	}

	// Generate response.
	rsp := &pb.CreateUserResponse{
		User: &pb.User{
			Username:         txResult.User.Username,
			FullName:         txResult.User.FullName,
			Email:            txResult.User.Email,
			PasswordChangeAt: timestamppb.New(txResult.User.PasswordChangeAt),
			CreateAt:         timestamppb.New(txResult.User.CreateAt),
		},
	}

	return rsp, nil
}

// validateCreateUserRequest validates all fields in CreateUserRequest.
func validateCreateUserRequest(req *pb.CreateUserRequest) (voilations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		voilations = append(voilations, fieldValidation("username", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		voilations = append(voilations, fieldValidation("password", err))
	}
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		voilations = append(voilations, fieldValidation("email", err))
	}
	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		voilations = append(voilations, fieldValidation("full_name", err))
	}
	return
}
