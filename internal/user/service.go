package user

type Service interface {
	// Login(ctx context.Context, username, password string)
}

type service struct {
}

func NewService() Service {
	return &service{}
}
