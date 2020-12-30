package main

import (
	"github.com/embersyndicate/support"
	"github.com/embersyndicate/support/internal/mongo"
)

type repositories struct {
	category support.CategoryRepository
	ticket   support.TicketRepository
	user     support.UserRepository
}

func initializeRepositories(basics *app) repositories {

	repos := repositories{}

	repos.category, err = mongo.NewCategoryRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize category repository")
	}

	basics.logger.Info("category repository initialized")

	repos.ticket, err = mongo.NewTicketRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize ticket repository")
	}

	basics.logger.Info("ticket repository initialized")

	repos.user, err = mongo.NewUserRepository(basics.db)
	if err != nil {
		basics.logger.WithError(err).Fatal("failed to initialize user repository")
	}

	basics.logger.Info("user repository initialized")

	return repos

}
