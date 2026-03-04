package repository

import (
	"database/sql"

	"github.com/captiv8io/go-microservices/services/discovery/creatorpricing/entity"
)

type NotificationSearchQuery struct {
	JobID       sql.NullString
	ListID      sql.NullInt64
	Status      entity.NotificationStatus
	SubmittedAt sql.NullTime
	CompletedAt sql.NullTime
	Skip        int32
	Take        int32
}
