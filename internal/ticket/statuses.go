package ticket

import (
	"context"
	"fmt"
	"time"

	"github.com/embersyndicate/support"
	"github.com/embersyndicate/support/pkg/middleware"
)

func (s *service) TicketStatus(ctx context.Context, id string) (*support.TicketStatus, error) {
	status, err := s.TicketRepository.TicketStatus(ctx, id)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to fetch status %s", id)
	}

	return status, nil
}

func (s *service) TicketStatuses(ctx context.Context, operators ...*support.Operator) ([]*support.TicketStatus, error) {

	statuses, err := s.TicketRepository.TicketStatuses(ctx, operators...)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to fetch statuses")
	}

	return statuses, nil

}

func (s *service) CreateTicketStatus(ctx context.Context, status *support.TicketStatus) (*support.TicketStatus, error) {

	err := status.ValidateAttributes()
	if err != nil {
		return nil, err
	}

	userID, err := middleware.GetUserObjectIDFromContext(ctx)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to retrieve user id from context")
	}

	now := time.Now()
	status.CreatedAt = now
	status.CreatedBy = userID
	status.UpdatedAt = now
	status.UpdatedBy = userID

	status, err = s.TicketRepository.CreateTicketStatus(ctx, status)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to create ticket status")
	}

	return status, err

}

func (s *service) UpdateTicketStatus(ctx context.Context, id string, status *support.TicketStatus) (*support.TicketStatus, error) {

	err := status.ValidateAttributes()
	if err != nil {
		return nil, err
	}

	userID, err := middleware.GetUserObjectIDFromContext(ctx)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to retrieve user id from context")
	}

	status.UpdatedAt = time.Now()
	status.UpdatedBy = userID

	status, err = s.TicketRepository.UpdateTicketStatus(ctx, id, status)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to update ticket status %s", id)

	}

	return status, err
}
