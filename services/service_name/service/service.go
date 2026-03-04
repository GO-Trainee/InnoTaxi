package service

type Service interface {
	// Define service methods here, for example:
	// GetUser(ctx context.Context, id string) (*serviceEntity.User, error)
	// CreateUser(ctx context.Context, user *serviceEntity.User) (*serviceEntity.User, error)
	// UpdateUser(ctx context.Context, id string, user *serviceEntity.User) (*serviceEntity.User, error)
}

type service struct {
	// Add any dependencies here, for example:
	// repo repository.Repository
	// gateway gateway.Gateway
}

func New( /*dependencies*/ ) Service {
	return &service{
		// repo: repo,
		// gateway: gateway,
	}
}
