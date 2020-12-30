package ticket

import (
	"context"
	"fmt"
	"time"

	"github.com/embersyndicate/support"
	"github.com/embersyndicate/support/internal"
	"github.com/embersyndicate/support/pkg/middleware"
)

func (s *service) TicketDefinition(ctx context.Context, id string) (*support.TicketDefinition, error) {

	definition, err := s.TicketRepository.TicketDefinition(ctx, id)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to fetch definition %s", id)
	}

	return definition, nil

}

func (s *service) TicketDefinitions(ctx context.Context, operators ...*support.Operator) ([]*support.TicketDefinition, error) {

	definitions, err := s.TicketRepository.TicketDefinitions(ctx, operators...)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to fetch ticket definitions")
	}

	return definitions, nil

}

func (s *service) CreateTicketDefinition(ctx context.Context, definition *support.TicketDefinition) (*support.TicketDefinition, error) {

	err := definition.ValidateAttributes()
	if err != nil {
		return nil, internal.NewInternalError(internal.LevelBad, err.Error())
	}

	userID, err := middleware.GetUserObjectIDFromContext(ctx)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to retrieve user id from context")
	}

	if len(definition.Fields) == 0 {
		return nil, internal.NewInternalError(internal.LevelBad, "field must have a length greater than or equal to 1, length of 0 detected")
	}

	for i, field := range definition.Fields {
		for j, ifield := range definition.Fields {
			if field.Hex() == ifield.Hex() && i != j {
				return nil, internal.NewInternalError(internal.LevelBad, "fields must be unique. tickets cannot have multiple fields with the same name")
			}
		}

		_, err := s.FieldDefinition(ctx, field.Hex())
		if err != nil {
			middleware.LogEntrySetError(ctx, err)
			return nil, internal.NewInternalError(internal.LevelBad, fmt.Sprintf("unable to resolve %s field id to valid field definition", field.Hex()))
		}

	}

	now := time.Now()
	definition.CreatedAt = now
	definition.CreatedBy = userID
	definition.UpdatedAt = now
	definition.UpdatedBy = userID

	definition, err = s.TicketRepository.CreateTicketDefinition(ctx, definition)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		if internal.IsUniqueConstrainViolation(err) {
			return nil, internal.NewInternalError(internal.LevelBad, "definition name must be unique")
		}
		return nil, internal.NewInternalError(internal.LevelInternal, "failed to create ticket definition")
	}

	return definition, nil

}

func (s *service) UpdateTicketDefinition(ctx context.Context, id string, definition *support.TicketDefinition) (*support.TicketDefinition, error) {

	err := definition.ValidateAttributes()
	if err != nil {
		return nil, err
	}

	currentDefinition, err := s.TicketDefinition(ctx, id)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, internal.NewInternalError(internal.LevelInternal, fmt.Sprintf("failed to fetch definition %s", id))
	}

	// Ensure that all the current field definitions exist in the updated definition
	for _, currentFieldDefinitionID := range currentDefinition.Fields {
		var exists bool
		for _, updatedFieldDefinitionID := range definition.Fields {
			if currentFieldDefinitionID.Hex() == updatedFieldDefinitionID.Hex() {
				exists = true
			}
		}
		if !exists {
			return nil, internal.NewInternalError(internal.LevelBad, "failed to updated definition. existing field definition missing from updated payload.")
		}
	}

	userID, err := middleware.GetUserObjectIDFromContext(ctx)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to retrieve user id from context")
	}

	now := time.Now()
	definition.UpdatedAt = now
	definition.UpdatedBy = userID

	definition, err = s.TicketRepository.UpdateTicketDefinition(ctx, id, definition)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to create definition")
	}

	return definition, nil

}
