package support

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryRepository interface {
	Category(ctx context.Context, id string) (*Category, error)
	Categories(ctx context.Context) ([]*Category, error)
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	UpdateCategory(ctx context.Context, id string, category *Category) (*Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

type Category struct {
	ID        primitive.ObjectID  `json:"id" bson:"_id"`
	ParentID  *primitive.ObjectID `json:"parentID" bson:"parentID"`
	Name      string              `json:"name" bson:"name"`
	Children  []*Category         `json:"children" bson:"-"`
	CreatedBy primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy *primitive.ObjectID `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt *time.Time          `json:"updatedAt" bson:"updatedAt"`
}
