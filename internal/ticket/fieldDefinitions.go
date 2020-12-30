package ticket

import (
	"context"
	"fmt"
	"time"

	"github.com/embersyndicate/support"
	"github.com/embersyndicate/support/pkg/middleware"
)

func (s *service) FieldDefinition(ctx context.Context, id string) (*support.FieldDefinition, error) {

	definition, err := s.TicketRepository.FieldDefinition(ctx, id)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to fetch definition %s", id)
	}

	return definition, nil

}

func (s *service) FieldDefinitions(ctx context.Context, operators ...*support.Operator) ([]*support.FieldDefinition, error) {

	definitions, err := s.TicketRepository.FieldDefinitions(ctx, operators...)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to fetch ticket definitions")
	}

	return definitions, nil

}

func (s *service) CreateFieldDefinition(ctx context.Context, definition *support.FieldDefinition) (*support.FieldDefinition, error) {

	err := definition.ValidateAttributes()
	if err != nil {
		return nil, err
	}

	userID, err := middleware.GetUserObjectIDFromContext(ctx)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to retrieve user id from context")
	}

	now := time.Now()
	definition.CreatedAt = now
	definition.CreatedBy = userID
	definition.UpdatedAt = now
	definition.UpdatedBy = userID

	definition, err = s.TicketRepository.CreateFieldDefinition(ctx, definition)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to create definition")
	}

	return definition, nil

}

func (s *service) UpdateFieldDefinition(ctx context.Context, id string, definition *support.FieldDefinition) (*support.FieldDefinition, error) {

	err := definition.ValidateAttributes()
	if err != nil {
		return nil, err
	}

	userID, err := middleware.GetUserObjectIDFromContext(ctx)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to retrieve user id from context")
	}

	now := time.Now()
	definition.UpdatedAt = now
	definition.UpdatedBy = userID

	definition, err = s.TicketRepository.UpdateFieldDefinition(ctx, id, definition)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to create definition")
	}

	return definition, nil

}
