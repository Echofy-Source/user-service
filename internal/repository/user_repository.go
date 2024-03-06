package repository

import (
	"context"

	"github.com/Echofy-Source/user-service/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	SearchUsers(ctx context.Context, username string) ([]*model.User, error)
}

type UserRepositoryImpl struct {
	collection *mongo.Collection
}

// CreateUser implements UserRepository.
func (u *UserRepositoryImpl) CreateUser(ctx context.Context, user *model.User) error {
	results, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.ID = results.InsertedID.(primitive.ObjectID)
	return nil
}

// GetUserByID implements UserRepository.
func (u *UserRepositoryImpl) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user model.User
	err = u.collection.FindOne(ctx, model.User{ID: objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername implements UserRepository.
func (u *UserRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := u.collection.FindOne(ctx, model.User{Username: username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SearchUsers implements UserRepository.
func (u *UserRepositoryImpl) SearchUsers(ctx context.Context, username string) ([]*model.User, error) {
	cursor, err := u.collection.Find(ctx, model.User{Username: username})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var users []*model.User
	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// UpdateUser implements UserRepository.
func (u *UserRepositoryImpl) UpdateUser(ctx context.Context, user *model.User) error {
	_, err := u.collection.ReplaceOne(ctx, model.User{ID: user.ID}, user)
	if err != nil {
		return err
	}
	return nil
}

func NewUserRepository(collection *mongo.Collection) *UserRepositoryImpl {
	return &UserRepositoryImpl{collection: collection}
}

// Ensure that UserRepositoryImpl implements UserRepository
var _ UserRepository = &UserRepositoryImpl{}
