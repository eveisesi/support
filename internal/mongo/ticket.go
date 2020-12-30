package mongo

import (
	"context"
	"fmt"

	"github.com/embersyndicate/support"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ticketRepository struct {
	tickets           *mongo.Collection
	ticketDefinitions *mongo.Collection
	ticketStatuses    *mongo.Collection
	fieldDefinitions  *mongo.Collection
}

func NewTicketRepository(d *mongo.Database) (support.TicketRepository, error) {

	t := d.Collection("tickets")
	td := d.Collection("ticketDefinitions")
	ts := d.Collection("ticketStatuses")
	tf := d.Collection("fieldDefinitions")

	_, err := td.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys: bson.M{
				"name": 1,
			},
			Options: &options.IndexOptions{
				Name:   newString("uniqueTicketDefinitionName"),
				Unique: newBool(true),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	_, err = ts.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys: bson.M{
				"name": 1,
			},
			Options: &options.IndexOptions{
				Name:   newString("uniqueTicketStatusName"),
				Unique: newBool(true),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	_, err = tf.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys: bson.M{
				"name": 1,
			},
			Options: &options.IndexOptions{
				Name:   newString("uniqueFieldDefinitionName"),
				Unique: newBool(true),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return &ticketRepository{
		tickets:           t,
		ticketDefinitions: td,
		ticketStatuses:    ts,
		fieldDefinitions:  tf,
	}, nil

}

func (r *ticketRepository) Ticket(ctx context.Context, id string) (*support.Ticket, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unable to cast %s to ObjectID", id)
	}

	tickets, err := r.Tickets(ctx, support.NewEqualOperator("_id", _id), support.NewLimitOperator(1))
	if err != nil {
		return nil, err
	}

	if len(tickets) == 0 {
		return nil, fmt.Errorf("category does not exist")
	}

	return tickets[0], nil
}

func (r *ticketRepository) Tickets(ctx context.Context, operators ...*support.Operator) ([]*support.Ticket, error) {
	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var tickets = make([]*support.Ticket, 0)
	result, err := r.tickets.Find(ctx, filters, options)
	if err != nil {
		return tickets, err
	}

	err = result.All(ctx, &tickets)

	return tickets, err
}

func (r *ticketRepository) CreateTicket(ctx context.Context, ticket *support.Ticket) (*support.Ticket, error) {

	result, err := r.tickets.InsertOne(ctx, ticket)
	if err != nil {
		return nil, err
	}

	ticket.ID = result.InsertedID.(primitive.ObjectID)

	return ticket, err

}

func (r *ticketRepository) UpdateTicket(ctx context.Context, id string, ticket *support.Ticket) (*support.Ticket, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unable to cast %s to ObjectID", id)
	}

	ticket.ID = _id

	update := primitive.D{primitive.E{Key: "$set", Value: ticket}}

	_, err = r.tickets.UpdateOne(ctx, primitive.D{primitive.E{Key: "_id", Value: _id}}, update)

	return ticket, err

}

func (r *ticketRepository) TicketDefinition(ctx context.Context, id string) (*support.TicketDefinition, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unable to cast %s to ObjectID", id)
	}

	ticketDefinitions, err := r.TicketDefinitions(ctx, support.NewEqualOperator("_id", _id), support.NewLimitOperator(1))
	if err != nil {
		return nil, err
	}

	if len(ticketDefinitions) == 0 {
		return nil, fmt.Errorf("category does not exist")
	}

	return ticketDefinitions[0], nil

}

func (r *ticketRepository) TicketDefinitions(ctx context.Context, operators ...*support.Operator) ([]*support.TicketDefinition, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var ticketDefinitions = make([]*support.TicketDefinition, 0)
	result, err := r.ticketDefinitions.Find(ctx, filters, options)
	if err != nil {
		return ticketDefinitions, err
	}

	err = result.All(ctx, &ticketDefinitions)

	return ticketDefinitions, err

}

func (r *ticketRepository) CreateTicketDefinition(ctx context.Context, ticketDefinition *support.TicketDefinition) (*support.TicketDefinition, error) {

	result, err := r.ticketDefinitions.InsertOne(ctx, ticketDefinition)
	if err != nil {
		return nil, err
	}

	ticketDefinition.ID = result.InsertedID.(primitive.ObjectID)

	return ticketDefinition, err

}

func (r *ticketRepository) UpdateTicketDefinition(ctx context.Context, id string, ticketDefinition *support.TicketDefinition) (*support.TicketDefinition, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unable to cast %s to ObjectID", id)
	}

	ticketDefinition.ID = _id

	update := primitive.D{primitive.E{Key: "$set", Value: ticketDefinition}}

	_, err = r.ticketDefinitions.UpdateOne(ctx, primitive.D{primitive.E{Key: "_id", Value: _id}}, update)

	return ticketDefinition, err

}

func (r *ticketRepository) TicketStatus(ctx context.Context, id string) (*support.TicketStatus, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unable to cast %s to ObjectID", id)
	}

	ticketStatuses, err := r.TicketStatuses(ctx, support.NewEqualOperator("_id", _id), support.NewLimitOperator(1))
	if err != nil {
		return nil, err
	}

	if len(ticketStatuses) == 0 {
		return nil, fmt.Errorf("category does not exist")
	}

	return ticketStatuses[0], nil

}

func (r *ticketRepository) TicketStatuses(ctx context.Context, operators ...*support.Operator) ([]*support.TicketStatus, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var ticketStatuses = make([]*support.TicketStatus, 0)
	result, err := r.ticketStatuses.Find(ctx, filters, options)
	if err != nil {
		return ticketStatuses, err
	}

	err = result.All(ctx, &ticketStatuses)

	return ticketStatuses, err

}

func (r *ticketRepository) CreateTicketStatus(ctx context.Context, ticketStatus *support.TicketStatus) (*support.TicketStatus, error) {

	result, err := r.ticketStatuses.InsertOne(ctx, ticketStatus)
	if err != nil {
		return nil, err
	}

	ticketStatus.ID = result.InsertedID.(primitive.ObjectID)

	return ticketStatus, err

}

func (r *ticketRepository) UpdateTicketStatus(ctx context.Context, id string, ticketStatus *support.TicketStatus) (*support.TicketStatus, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unable to cast %s to ObjectID", id)
	}

	ticketStatus.ID = _id

	update := primitive.D{primitive.E{Key: "$set", Value: ticketStatus}}

	_, err = r.ticketStatuses.UpdateOne(ctx, primitive.D{primitive.E{Key: "_id", Value: _id}}, update)

	return ticketStatus, err

}

func (r *ticketRepository) FieldDefinition(ctx context.Context, id string) (*support.FieldDefinition, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unable to cast %s to ObjectID", id)
	}

	definitions, err := r.FieldDefinitions(ctx, support.NewEqualOperator("_id", _id), support.NewLimitOperator(1))
	if err != nil {
		return nil, err
	}

	if len(definitions) == 0 {
		return nil, fmt.Errorf("category does not exist")
	}

	return definitions[0], nil

}

func (r *ticketRepository) FieldDefinitions(ctx context.Context, operators ...*support.Operator) ([]*support.FieldDefinition, error) {

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	var definitions = make([]*support.FieldDefinition, 0)
	result, err := r.fieldDefinitions.Find(ctx, filters, options)
	if err != nil {
		return definitions, err
	}

	err = result.All(ctx, &definitions)

	return definitions, err

}

func (r *ticketRepository) CreateFieldDefinition(ctx context.Context, definition *support.FieldDefinition) (*support.FieldDefinition, error) {

	result, err := r.fieldDefinitions.InsertOne(ctx, definition)
	if err != nil {
		return nil, err
	}

	definition.ID = result.InsertedID.(primitive.ObjectID)

	return definition, err

}

func (r *ticketRepository) UpdateFieldDefinition(ctx context.Context, id string, definition *support.FieldDefinition) (*support.FieldDefinition, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("unable to cast %s to ObjectID", id)
	}

	definition.ID = _id

	update := primitive.D{primitive.E{Key: "$set", Value: definition}}

	_, err = r.fieldDefinitions.UpdateOne(ctx, primitive.D{primitive.E{Key: "_id", Value: _id}}, update)

	return definition, err

}
