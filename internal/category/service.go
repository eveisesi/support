package category

import (
	"context"
	"fmt"
	"time"

	"github.com/embersyndicate/support"

	"github.com/embersyndicate/support/pkg/middleware"
)

type Service interface {
	support.CategoryRepository
}

type service struct {
	// cache  support.CategoryRepository
	support.CategoryRepository
}

func New(category support.CategoryRepository) Service {

	s := &service{
		CategoryRepository: category,
	}

	return s
}

func (s *service) Category(ctx context.Context, id string) (*support.Category, error) {

	category, err := s.CategoryRepository.Category(ctx, id)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to fetch category %s", id)
	}

	return category, nil

}

func (s *service) Categories(ctx context.Context, operators ...*support.Operator) ([]*support.Category, error) {

	categories, err := s.CategoryRepository.Categories(ctx, operators...)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to fetch categories")
	}

	return categories, nil

}

func (s *service) CreateCategory(ctx context.Context, category *support.Category) (*support.Category, error) {

	err := category.VerifyAttributes()
	if err != nil {
		return nil, err
	}

	userID, err := middleware.GetUserObjectIDFromContext(ctx)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to retrieve user id from context")
	}

	now := time.Now()
	category.CreatedAt = now
	category.CreatedBy = userID
	category.UpdatedAt = now
	category.UpdatedBy = userID

	category, err = s.CategoryRepository.CreateCategory(ctx, category)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, err

}

func (s *service) UpdateCategory(ctx context.Context, id string, category *support.Category) (*support.Category, error) {

	err := category.VerifyAttributes()
	if err != nil {
		return nil, err
	}

	userID, err := middleware.GetUserObjectIDFromContext(ctx)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	category.UpdatedAt = time.Now()
	category.UpdatedBy = userID

	category, err = s.CategoryRepository.UpdateCategory(ctx, id, category)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to update category %s", id)
	}

	return category, err

}
