package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/util"
	"github.com/techschool/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// Validate update user request.
	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	// Update user follows request.
	arg := db.UpdateUserParams{
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
		Username: req.GetUsername(),
	}

	if req.Password != nil {
		hashedPassword, err := util.HashedPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}
		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
		// Update field password_change_at.
		arg.PasswordChangeAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Unimplemented, "failed to update user: %s", err)
	}

	// Generate response.
	rsp := &pb.UpdateUserResponse{
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

// validateUpdateUserRequest validates all fields in UpdateUserRequest.
func validateUpdateUserRequest(req *pb.UpdateUserRequest) (voilations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		voilations = append(voilations, fieldValidation("username", err))
	}
	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			voilations = append(voilations, fieldValidation("password", err))
		}
	}
	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			voilations = append(voilations, fieldValidation("email", err))
		}
	}
	if req.FullName != nil {
		if err := val.ValidateFullName(req.GetFullName()); err != nil {
			voilations = append(voilations, fieldValidation("full_name", err))
		}
	}
	return
}
