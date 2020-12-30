package ticket

import (
	"github.com/embersyndicate/support"
)

type Service interface {
	support.TicketRepository
}

type service struct {
	support.TicketRepository
}

func New(ticket support.TicketRepository) Service {
	return &service{
		TicketRepository: ticket,
	}
}
