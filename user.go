package support

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	User(ctx context.Context, id string) (*User, error)
	Users(ctx context.Context) ([]*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, id string, user *User) (*User, error)
	DeleteUser(ctx context.Context, id string) error
}

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"password" bson:"password"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
