package admin_users_api

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/stretchr/testify/require"
)

// fakeAdminUsersService is a minimal in-memory service.AdminUsersService used to test
// AdminUsersImpl's handlers without a real database — following the fakeAuthService convention
// in internal/transport/auth_api/auth_impl_test.go (a plain struct with configurable
// result/err fields per method, not an embedded-interface trick).
type fakeAdminUsersService struct {
	createUserResult domain.User
	createUserErr    error

	changePasswordErr        error
	changePasswordCalledWith struct {
		userUuid    uuid.UUID
		newPassword string
	}
}

func (f *fakeAdminUsersService) ListUsers(
	_ context.Context, _ domain.ListUsersReq,
) ([]domain.User, int64, error) {
	panic("not implemented")
}

func (f *fakeAdminUsersService) GetUser(_ context.Context, _ uuid.UUID) (domain.UserDetails, error) {
	panic("not implemented")
}

func (f *fakeAdminUsersService) GetUserSessions(_ context.Context, _ uuid.UUID) ([]domain.Session, error) {
	panic("not implemented")
}

func (f *fakeAdminUsersService) CreateUser(_ context.Context, _, _ string) (domain.User, error) {
	if f.createUserErr != nil {
		return domain.User{}, f.createUserErr
	}

	return f.createUserResult, nil
}

func (f *fakeAdminUsersService) ChangePassword(_ context.Context, userUuid uuid.UUID, newPassword string) error {
	f.changePasswordCalledWith.userUuid = userUuid
	f.changePasswordCalledWith.newPassword = newPassword

	return f.changePasswordErr
}

var _ service.AdminUsersService = (*fakeAdminUsersService)(nil)

func TestAdminUsersImpl_CreateArtelUser(t *testing.T) {
	userUuid := uuid.New()
	svc := &fakeAdminUsersService{
		createUserResult: domain.User{
			Uuid:     userUuid,
			Email:    "new-user@example.com",
			Username: "new-user",
		},
	}
	impl := New(svc, false)

	req := &artel_api.CreateArtelUser_Request{Email: "new-user@example.com", Password: "s3cret-p4ss"}

	resp, err := impl.CreateArtelUser(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp.User)
	require.Equal(t, userUuid.String(), resp.User.UserId)
	require.Equal(t, "new-user@example.com", resp.User.Email)
	require.Equal(t, "new-user", resp.User.Username)
}

func TestAdminUsersImpl_CreateArtelUser_PropagatesError(t *testing.T) {
	svcErr := errors.New("email already registered")
	svc := &fakeAdminUsersService{createUserErr: svcErr}
	impl := New(svc, false)

	req := &artel_api.CreateArtelUser_Request{Email: "dup@example.com", Password: "s3cret-p4ss"}

	resp, err := impl.CreateArtelUser(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAdminUsersImpl_ChangeArtelUserPassword(t *testing.T) {
	userUuid := uuid.New()
	svc := &fakeAdminUsersService{}
	impl := New(svc, false)

	requestUuidStr := userUuid.String()
	req := &artel_api.ChangeArtelUserPassword_Request{UserId: requestUuidStr, NewPassword: "n3w-p4ss"}

	resp, err := impl.ChangeArtelUserPassword(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, userUuid, svc.changePasswordCalledWith.userUuid)
	require.Equal(t, "n3w-p4ss", svc.changePasswordCalledWith.newPassword)
}

func TestAdminUsersImpl_ChangeArtelUserPassword_InvalidUserId(t *testing.T) {
	svc := &fakeAdminUsersService{}
	impl := New(svc, false)

	req := &artel_api.ChangeArtelUserPassword_Request{UserId: "not-a-uuid", NewPassword: "n3w-p4ss"}

	resp, err := impl.ChangeArtelUserPassword(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAdminUsersImpl_ChangeArtelUserPassword_PropagatesError(t *testing.T) {
	svcErr := errors.New("password too weak")
	svc := &fakeAdminUsersService{changePasswordErr: svcErr}
	impl := New(svc, false)

	req := &artel_api.ChangeArtelUserPassword_Request{UserId: uuid.New().String(), NewPassword: "weak"}

	resp, err := impl.ChangeArtelUserPassword(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, resp)
}
