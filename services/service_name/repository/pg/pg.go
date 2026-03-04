package pg

type PgRepository interface {
	PgAtomicRepository
	// Define your repository methods here, for example:
	// Create(ctx context.Context, tx *sqlxBase.Tx, entity *repositoryEnity.User) (*repositoryEnity.User, error)
	// Update(ctx context.Context, tx *sqlxBase.Tx, entity *repositoryEnity.User) (*repositoryEnity.User, error)
	// FetchById(ctx context.Context, id int64) (*repositoryEnity.User, error)
}

type pgRepository struct {
	// Add any dependencies here, for example:
	// db *sqlx.DB
}

func NewPgRepository( /* Add any dependencies here, for example: db *sqlx.DB */ ) PgRepository {
	return &pgRepository{
		// Initialize your dependencies here, for example: db: db,
	}
}
