package repository

import (
	"awesomeProject/services/service_name/repository/mongo"
	"awesomeProject/services/service_name/repository/pg"
	"awesomeProject/services/service_name/repository/redis"
)

type Repository struct {
	Redis redis.RedisRepository
	Mongo mongo.MongoRepository
	Pg    pg.PgRepository
}

func New( /*redisRepo RedisRepository, mongoRepo MongoRepository, pgRepo PgRepository*/ ) Repository {
	return Repository{
		// Redis: redisRepo,
		// Mongo: mongoRepo,
		// Pg:   pgRepo,
	}
}
