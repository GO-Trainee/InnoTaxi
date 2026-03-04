package mongo

type MongoClient struct {
	//
}

// New creates the Mongo client and immediately runs migrations.
// Migrations must complete successfully before the app proceeds —
// if they fail, New must return an error and the service must not start.
func New() *MongoClient {
	// TODO: connect to DB
	// TODO: call runMigrations() here — do NOT defer it to later in run.go
	return &MongoClient{}
}

func runMigrations() {}
