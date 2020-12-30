package support

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TicketRepository interface {
	ticketRepository
	ticketDefinitionRepository
	ticketStatusRepository
	fieldDefinitionRepository
}

type ticketRepository interface {
	Ticket(ctx context.Context, id string) (*Ticket, error)
	Tickets(ctx context.Context, operators ...*Operator) ([]*Ticket, error)
	CreateTicket(ctx context.Context, ticket *Ticket) (*Ticket, error)
	UpdateTicket(ctx context.Context, id string, ticket *Ticket) (*Ticket, error)
}

type ticketDefinitionRepository interface {
	TicketDefinition(ctx context.Context, id string) (*TicketDefinition, error)
	TicketDefinitions(ctx context.Context, operators ...*Operator) ([]*TicketDefinition, error)
	CreateTicketDefinition(ctx context.Context, ticket *TicketDefinition) (*TicketDefinition, error)
	UpdateTicketDefinition(ctx context.Context, id string, ticket *TicketDefinition) (*TicketDefinition, error)
}

type ticketStatusRepository interface {
	TicketStatus(ctx context.Context, id string) (*TicketStatus, error)
	TicketStatuses(ctx context.Context, operators ...*Operator) ([]*TicketStatus, error)
	CreateTicketStatus(ctx context.Context, ticket *TicketStatus) (*TicketStatus, error)
	UpdateTicketStatus(ctx context.Context, id string, ticket *TicketStatus) (*TicketStatus, error)
}

type fieldDefinitionRepository interface {
	FieldDefinition(ctx context.Context, id string) (*FieldDefinition, error)
	FieldDefinitions(ctx context.Context, operators ...*Operator) ([]*FieldDefinition, error)
	CreateFieldDefinition(ctx context.Context, definition *FieldDefinition) (*FieldDefinition, error)
	UpdateFieldDefinition(ctx context.Context, id string, ticket *FieldDefinition) (*FieldDefinition, error)
}

// Ticket is represents a ticket that has been submitted to this system.
// The TypeID dictates the fields that should be attached to it, who created it, and who the ticket is assigned to.
type Ticket struct {
	ID           primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	SubmittedBy  primitive.ObjectID  `json:"submittedBy" bson:"submittedBy"`
	AssignedTo   *primitive.ObjectID `json:"assignedTo,omitempty" bson:"assignedTo,omitempty"`
	StatusID     primitive.ObjectID  `json:"statusID" bson:"statusID"`
	DefinitionID primitive.ObjectID  `json:"definitionID" bson:"definitionID"`
	CategoryID   primitive.ObjectID  `json:"categoryID" bson:"categoryID"`
	Fields       []*FieldValue       `json:"fields" bson:"fields"`
	CreatedAt    time.Time           `json:"createdAt" bson:"createdAt"`
	UpdateAt     *time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

// TicketType represents a type of ticket and the fields that the ticket has
type TicketDefinition struct {
	ID         primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name       string               `json:"name" bson:"name"`
	Fields     []primitive.ObjectID `json:"fields" bson:"fields"`
	Disabled   bool                 `json:"disabled" bson:"disabled"`
	DisabledBy *primitive.ObjectID  `json:"disabledBy,omitempty" bson:"disabledBy,omitempty"`
	DisabledAt *time.Time           `json:"disabledAt,omitempty" bson:"disabledAt,omitempty"`
	CreatedBy  primitive.ObjectID   `json:"createdBy" bson:"createdBy"`
	CreatedAt  time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedBy  primitive.ObjectID   `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	UpdatedAt  time.Time            `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func (o *TicketDefinition) ValidateAttributes() error {

	if o.Name == "" {
		return fmt.Errorf("name is required, received empty value")
	}

	if len(o.Fields) == 0 {
		return fmt.Errorf("fields is required, received empty array")
	}

	return nil

}

type TicketStatus struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`

	// If tickets are placed into a status that is locked,
	// the status of that ticket cannot be updated without a TBD override.
	// Thought process is once a ticket has been closed it can't be reopened.
	Locked bool `json:"locked" bson:"locked"`

	CreatedBy primitive.ObjectID `json:"createdBy" bson:"createdBy"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedBy primitive.ObjectID `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func (o *TicketStatus) ValidateAttributes() error {

	if o.Name == "" {
		return fmt.Errorf("name is required, received empty value")
	}

	return nil

}

type FieldKind string

const (
	FieldString  FieldKind = "string"
	FieldNumber  FieldKind = "number"
	FieldBoolean FieldKind = "boolean"
	FieldList    FieldKind = "list"
)

type Kinds []FieldKind

var AllKinds = Kinds{
	FieldString, FieldNumber,
	FieldBoolean, FieldList,
}

func (a Kinds) Slice() []string {
	out := make([]string, len(a))
	for i, v := range a {
		out[i] = v.String()
	}
	return out
}

func (f FieldKind) Valid() bool {
	for _, v := range AllKinds {
		if v == f {
			return true
		}
	}

	return false
}

func (f FieldKind) String() string {
	return string(f)
}

// Field represents a field that is apart of a ticket
type FieldDefinition struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Required    bool               `json:"required" bson:"required"`
	Hidden      bool               `json:"hidden" bson:"hidden"`
	Hash        bool               `json:"hash" bson:"hash"`
	Kind        FieldKind          `json:"kind" bson:"kind"`
	Options     []interface{}      `json:"options,omitempty" bson:"options,omitempty"`
	Disabled    bool               `json:"disabled" bson:"disabled"`
	DisabledBy  primitive.ObjectID `json:"disabledBy,omitempty" bson:"disabledBy,omitempty"`
	DisabledAt  time.Time          `json:"disabledAt,omitempty" bson:"disabledAt,omitempty"`
	CreatedBy   primitive.ObjectID `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedBy   primitive.ObjectID `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func (o *FieldDefinition) ValidateAttributes() error {

	if o.Name == "" {
		return fmt.Errorf("name is required, received empty value")
	}

	if o.Description == "" {
		return fmt.Errorf("description is required, received empty value")
	}

	if o.Kind == "" {
		return fmt.Errorf("kind is required, received empty value")
	}

	if !o.Kind.Valid() {
		return fmt.Errorf("invalid value for kind provided, got %s, exported on of %s", o.Kind, strings.Join(AllKinds.Slice(), ", "))
	}

	if o.Kind == FieldList && len(o.Options) == 0 {
		return fmt.Errorf("options cannot be empty with kind is %s", FieldList)
	}

	return nil

}

type FieldValue struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Value interface{}        `json:"value" bson:"value"`
}
