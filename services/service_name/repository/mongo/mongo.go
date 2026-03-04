package mongo

type MongoRepository interface {
	// Define your repository methods here, for example:
	// Create(ctx context.Context, entity *repositoryEnity.User) (*repositoryEnity.User, error)
	// Update(ctx context.Context, entity *repositoryEnity.User) (*repositoryEnity.User, error)
	// FetchById(ctx context.Context, id string) (*repositoryEnity.User, error)
}

type mongoRepository struct {
	// Add any dependencies here, for example:
	// client *mongo.Client
}

func NewMongoRepository( /* Add any dependencies here, for example: client *mongo.Client */ ) MongoRepository {
	return &mongoRepository{
		// Initialize your dependencies here, for example: client: client,
	}
}
