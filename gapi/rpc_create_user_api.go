package gapi

import (
	"context"

	"github.com/lib/pq"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/util"
	"github.com/techschool/simplebank/val"
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

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}
	user, err := server.store.CreateUser(ctx, arg)

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
			Username:         user.Username,
			FullName:         user.FullName,
			Email:            user.Email,
			PasswordChangeAt: timestamppb.New(user.PasswordChangeAt),
			CreateAt:         timestamppb.New(user.CreateAt),
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
