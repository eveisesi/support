package ticket

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/embersyndicate/support/internal"

	"github.com/embersyndicate/support"
	"github.com/embersyndicate/support/pkg/middleware"
)

func (s *service) Ticket(ctx context.Context, id string) (*support.Ticket, error) {
	ticket, err := s.TicketRepository.Ticket(ctx, id)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to fetch ticket %s", id)
	}

	return ticket, nil
}

func (s *service) Tickets(ctx context.Context, operators ...*support.Operator) ([]*support.Ticket, error) {

	tickets, err := s.TicketRepository.Tickets(ctx, operators...)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to fetch tickets")
	}

	return tickets, nil

}

func (s *service) CreateTicket(ctx context.Context, ticket *support.Ticket) (*support.Ticket, error) {

	definition, err := s.TicketDefinition(ctx, ticket.DefinitionID.Hex())
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, internal.NewInternalError(internal.LevelBad, fmt.Sprintf("unknown definition %s", ticket.DefinitionID.Hex()))
	}

	fields, err := s.FieldDefinitions(ctx, support.NewInOperator("id", definition.Fields))
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, internal.NewInternalError(internal.LevelInternal, "failed to retrieve field definitions for specified ticket definition")
	}

	spew.Dump(fields)

	return nil, nil

}

func (s *service) UpdateTicket(ctx context.Context, id string, ticket *support.Ticket) (*support.Ticket, error) {

	panic("UpdateTicket not implemented")

}
