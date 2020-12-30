package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/embersyndicate/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	users *mongo.Collection
}

func NewUserRepository(d *mongo.Database) (support.UserRepository, error) {
	c := d.Collection("users")

	// Add Indexes if needed

	return &userRepository{
		users: c,
	}, nil
}

func (r *userRepository) User(ctx context.Context, id string) (*support.User, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unable to cast %s to ObjectID", id)
	}

	users, err := r.Users(ctx, support.NewEqualOperator("_id", _id), support.NewLimitOperator(1))
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user does not exist")
	}

	return users[0], nil

}

func (r *userRepository) Users(ctx context.Context, operators ...*support.Operator) ([]*support.User, error) {
	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var users = make([]*support.User, 0)
	result, err := r.users.Find(ctx, filters, options)
	if err != nil {
		return users, err
	}

	err = result.All(ctx, &users)

	return users, err
}

func (r *userRepository) CreateUser(ctx context.Context, user *support.User) (*support.User, error) {

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := r.users.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return user, err

}

func (r *userRepository) UpdateUser(ctx context.Context, id string, user *support.User) (*support.User, error) {
	panic("UpdateUser not implemented")
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	panic("DeleteUser not implemented")
}
