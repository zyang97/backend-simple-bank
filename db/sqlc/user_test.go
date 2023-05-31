package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomUser(t *testing.T) User {
	hashPassword, err := util.HashedPassword(util.RandString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandOwner(),
		HashedPassword: hashPassword,
		FullName:       util.RandOwner(),
		Email:          util.RandEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangeAt.IsZero())
	require.NotZero(t, user.CreateAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangeAt, user2.PasswordChangeAt, time.Second)
	require.WithinDuration(t, user1.CreateAt, user2.CreateAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	user := createRandomUser(t)
	fullName := util.RandOwner()
	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		FullName: sql.NullString{
			String: fullName,
			Valid:  true,
		},
		Username: user.Username,
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.Equal(t, newUser.FullName, fullName)
	require.Equal(t, user.Email, newUser.Email)
	require.Equal(t, user.HashedPassword, newUser.HashedPassword)

}

func TestUpdateUserOnlyEmail(t *testing.T) {
	user := createRandomUser(t)
	email := util.RandEmail()
	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Email: sql.NullString{
			String: email,
			Valid:  true,
		},
		Username: user.Username,
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.Equal(t, user.FullName, newUser.FullName)
	require.Equal(t, newUser.Email, email)
	require.Equal(t, user.HashedPassword, newUser.HashedPassword)

}

func TestUpdateUserOnlyPassword(t *testing.T) {
	user := createRandomUser(t)
	password := util.RandString(10)
	hashedPassword, err := util.HashedPassword(password)
	require.NoError(t, err)

	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
		Username: user.Username,
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.Equal(t, user.FullName, newUser.FullName)
	require.Equal(t, user.Email, newUser.Email)
	require.Equal(t, newUser.HashedPassword, hashedPassword)
}

func TestUpdateUserAllFields(t *testing.T) {
	user := createRandomUser(t)
	password := util.RandString(10)
	hashedPassword, err := util.HashedPassword(password)
	fullName := util.RandOwner()
	email := util.RandEmail()
	require.NoError(t, err)

	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: fullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: email,
			Valid:  true,
		},
		Username: user.Username,
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.Equal(t, newUser.FullName, fullName)
	require.Equal(t, newUser.Email, email)
	require.Equal(t, newUser.HashedPassword, hashedPassword)
}
