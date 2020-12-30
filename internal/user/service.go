package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/embersyndicate/support"
	"github.com/embersyndicate/support/internal/key"
	"github.com/embersyndicate/support/internal/token"
	"github.com/embersyndicate/support/pkg/middleware"
	"github.com/hesahesa/pwdbro"
	"github.com/hesahesa/pwdbro/checker"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(ctx context.Context, user *support.User) ([]byte, error)
	Register(ctx context.Context, user *support.User) (*support.User, error)
}

type service struct {
	pwdbro *pwdbro.PwdBro
	client *http.Client

	key   key.Service
	token token.Service

	userStore support.UserRepository
	// userCache support.UserRepository
}

func New(client *http.Client, key key.Service, token token.Service, user support.UserRepository) Service {

	pwd := pwdbro.NewEmptyPwdBro()

	// AddChecker always returns nil for err
	_ = pwd.AddChecker(&checker.Pwnedpasswords{
		HTTPClient: client,
	})

	s := &service{
		pwdbro: pwd,
		client: client,

		key:       key,
		token:     token,
		userStore: user,
	}

	return s
}

func (s *service) Login(ctx context.Context, user *support.User) ([]byte, error) {

	err := user.VerifyLoginAttributes()
	if err != nil {
		return nil, err
	}

	users, err := s.userStore.Users(
		ctx,
		support.NewOrOperator(
			support.NewEqualOperator(support.UserUsername, user.Username),
			support.NewEqualOperator(support.UserEmail, user.Email),
		),
		support.NewLimitOperator(1),
	)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to query for user")
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("username not found")
	}

	local := users[0]

	err = bcrypt.CompareHashAndPassword([]byte(local.Password), []byte(user.Password))
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("username/password combination is invalid")
	}

	key, err := s.token.BuildAndSignUserKey(ctx, local)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to generate token")
	}

	return key, nil

}

func (s *service) Register(ctx context.Context, user *support.User) (*support.User, error) {

	if err := user.VerifyRegisterAttributes(); err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, err
	}

	// Confirm that the user does not exist
	users, err := s.userStore.Users(ctx, support.NewOrOperator(support.NewEqualOperator(support.UserUsername, user.Username), support.NewEqualOperator(support.UserEmail, user.Email)))
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to query users for username")
	}

	if len(users) > 0 {
		return nil, fmt.Errorf("username is not unique")
	}

	err = s.checkPassword(ctx, user.Password)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, err
	}

	// If the username is unique and the password is not compromised or weak, lets replace the plain text password that was passed to us
	// with a hashed password
	user.Password, err = hashAndSaltPassword(ctx, user.Password)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, err
	}

	user, err = s.userStore.CreateUser(ctx, user)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return nil, fmt.Errorf("failed to register user")
	}

	// Sanitize the users password so it is not output upstream
	user.Password = ""

	return user, err
}

func hashAndSaltPassword(ctx context.Context, p string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return "", fmt.Errorf("failed to generate password hash")
	}

	return string(hash), nil

}

func (s *service) checkPassword(ctx context.Context, password string) error {

	if len(password) < 12 {
		err := fmt.Errorf("passwords must be atleast 12 chars long")
		middleware.LogEntrySetError(ctx, err)
		return err
	}

	// Check password strength
	statuses, err := s.pwdbro.RunChecks(password)
	if err != nil {
		middleware.LogEntrySetError(ctx, err)
		return fmt.Errorf("failed to validate password")
	}

	// Status is a slice of statuses whose length is equal to the number of checkers in teh instance of pwdbro.
	// Since we only have one checker registered, lets grab the entry at index 0
	status := statuses[0]
	if !status.Safe {
		middleware.LogEntrySetError(ctx, err)
		return fmt.Errorf("invalid or weak password detected")
	}

	return nil
}
