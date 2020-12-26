package support

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TicketRepository interface {
	ticketRepository
	ticketTypeRepository
}

type ticketRepository interface {
	Ticket(ctx context.Context, id string) (*Ticket, error)
	Tickets(ctx context.Context) ([]*Ticket, error)
	CreateTicket(ctx context.Context, ticket *Ticket) (*Ticket, error)
	UpdateTicket(ctx context.Context, id string, ticket *Ticket) (*Ticket, error)
}

type ticketTypeRepository interface {
	TicketType(ctx context.Context, id string) (*TicketType, error)
	TicketTypes(ctx context.Context) ([]*TicketType, error)
	CreateTicketType(ctx context.Context, ticket *TicketType) (*TicketType, error)
	UpdateTicketType(ctx context.Context, id string, ticket *TicketType) (*TicketType, error)
}

// Ticket is represents a ticket that has been submitted to this system.
// The TypeID dictates the fields that should be attached to it, who created it, and who the ticket is assigned to.
type Ticket struct {
	ID          primitive.ObjectID  `json:"id" bson:"_id"`
	SubmittedBy *primitive.ObjectID `json:"submittedBy" bson:"submittedBy"`
	AssignedTo  *primitive.ObjectID `json:"assignedTo" bson:"assignedTo"`
	StatusID    primitive.ObjectID  `json:"statusID" bson:"statusID"`
	TypeID      primitive.ObjectID  `json:"typeID" bson:"typeID"`
	TicketID    primitive.ObjectID  `json:"categoryID" bson:"categoryID"`
	Fields      []*Field            `json:"fields" bson:"fields"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdateAt    *time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

// TicketType represents a type of ticket and the fields that the ticket has
type TicketType struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Fields    []*Field           `json:"fields" bson:"fields"`
	CreatedBy primitive.ObjectID `json:"createdBy" bson:"createdBy"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedBy primitive.ObjectID `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
