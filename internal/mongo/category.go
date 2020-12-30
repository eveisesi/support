package mongo

import (
	"context"
	"fmt"

	"github.com/embersyndicate/support"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type categoryRepository struct {
	categories *mongo.Collection
}

func NewCategoryRepository(d *mongo.Database) (support.CategoryRepository, error) {

	c := d.Collection("categories")

	return &categoryRepository{
		categories: c,
	}, nil

}

func (r *categoryRepository) Category(ctx context.Context, id string) (*support.Category, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unable to cast %s to ObjectID", id)
	}

	categories, err := r.Categories(ctx, support.NewEqualOperator("_id", _id), support.NewLimitOperator(1))
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, fmt.Errorf("category does not exist")
	}

	return categories[0], nil
}

func (r *categoryRepository) Categories(ctx context.Context, operators ...*support.Operator) ([]*support.Category, error) {
	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var categories = make([]*support.Category, 0)
	result, err := r.categories.Find(ctx, filters, options)
	if err != nil {
		return categories, err
	}

	err = result.All(ctx, &categories)

	return categories, err
}

func (r *categoryRepository) CreateCategory(ctx context.Context, category *support.Category) (*support.Category, error) {

	result, err := r.categories.InsertOne(ctx, category)
	if err != nil {
		return nil, err
	}

	category.ID = result.InsertedID.(primitive.ObjectID)

	return category, err
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, id string, category *support.Category) (*support.Category, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unable to cast %s to ObjectID", id)
	}

	category.ID = _id

	update := primitive.D{primitive.E{Key: "$set", Value: category}}

	_, err = r.categories.UpdateOne(ctx, primitive.D{primitive.E{Key: "_id", Value: _id}}, update)

	return category, err
}
