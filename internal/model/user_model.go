package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Username      string             `bson:"username"`
	PasswordHash  string             `bson:"password_hash"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
	SignedPreKey  SignedPreKey       `bson:"signed_pre_key"`
	OneTimePreKey OneTimePreKey      `bson:"one_time_pre_key"`
}

type SignedPreKey struct {
	Key       string `bson:"key"`
	Signature string `bson:"signature"`
}

type OneTimePreKey struct {
	Key string `bson:"key"`
}

type CreateUserRequest struct {
	Username      string        `json:"username"`
	Password      string        `json:"password"`
	PublicKey     string        `json:"public_key"`
	SignedPreKey  SignedPreKey  `json:"signed_pre_key"`
	OneTimePreKey OneTimePreKey `json:"one_time_pre_key"`
}
