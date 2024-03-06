package service

import (
	"context"
	"time"

	"github.com/Echofy-Source/user-service/internal/model"
	"github.com/Echofy-Source/user-service/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, req *model.CreateUserRequest) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	SearchUsers(ctx context.Context, username string) ([]*model.User, error)
}

type UserServiceImpl struct {
	repo   repository.UserRepository
	crypto CryptoService
}

// CreateUser implements UserService.
func (u *UserServiceImpl) CreateUser(ctx context.Context, req *model.CreateUserRequest) error {
	hashedPassword, err := u.crypto.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &model.User{
		Username:      req.Username,
		PasswordHash:  hashedPassword,
		SignedPreKey:  req.SignedPreKey,
		OneTimePreKey: req.OneTimePreKey,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return u.repo.CreateUser(ctx, user)
}

// GetUserByID implements UserService.
func (u *UserServiceImpl) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return u.repo.GetUserByID(ctx, id)
}

// GetUserByUsername implements UserService.
func (u *UserServiceImpl) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return u.repo.GetUserByUsername(ctx, username)
}

// SearchUsers implements UserService.
func (u *UserServiceImpl) SearchUsers(ctx context.Context, username string) ([]*model.User, error) {
	return u.repo.SearchUsers(ctx, username)
}

// UpdateUser implements UserService.
func (u *UserServiceImpl) UpdateUser(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()
	return u.repo.UpdateUser(ctx, user)
}

// Ensure UserServiceImpl implements UserService.
var _ UserService = &UserServiceImpl{}
