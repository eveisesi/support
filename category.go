package support

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryRepository interface {
	Category(ctx context.Context, id string) (*Category, error)
	Categories(ctx context.Context, operators ...*Operator) ([]*Category, error)
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	UpdateCategory(ctx context.Context, id string, category *Category) (*Category, error)
}

type Category struct {
	ID        primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	ParentID  *primitive.ObjectID `json:"parentID" bson:"parentID"`
	Name      string              `json:"name" bson:"name"`
	CreatedBy primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt time.Time           `json:"updatedAt" bson:"updatedAt"`
}

func (o *Category) VerifyAttributes() error {
	if o.Name == "" {
		return fmt.Errorf("name is required, received empty value")
	}

	return nil
}
