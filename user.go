package support

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	User(ctx context.Context, id string) (*User, error)
	Users(ctx context.Context, operators ...*Operator) ([]*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, id string, user *User) (*User, error)
	DeleteUser(ctx context.Context, id string) error
}

// The following is a const list of the column name
// for each user struct filed that we tell mongo to use
const (
	UserUsername = "username"
	UserEmail    = "email"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName string             `json:"first_name" bson:"first_name"`
	LastName  string             `json:"last_name" bson:"last_name"`
	Email     string             `json:"email" bson:"email"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"password,omitempty" bson:"password"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func (o *User) VerifyLoginAttributes() error {

	if o.Username == "" {
		return fmt.Errorf("username required, received empty value")
	}

	if o.Password == "" {
		return fmt.Errorf("password required, received empty value")
	}

	return nil

}

func (o *User) VerifyRegisterAttributes() error {

	if o.FirstName == "" {
		return fmt.Errorf("first name required, received empty value")
	}

	if o.LastName == "" {
		return fmt.Errorf("last name required, received empty value")
	}

	if o.Email == "" {
		return fmt.Errorf("email address required, received empty value")
	}

	if o.Username == "" {
		return fmt.Errorf("username required, received empty value")
	}

	if o.Password == "" {
		return fmt.Errorf("password required, received empty value")
	}

	return nil

}
