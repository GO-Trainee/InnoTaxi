package repository

import (
	"database/sql"
)

type NotificationSearchQuery struct {
	JobID       sql.NullString
	ListID      sql.NullInt64
	Status      types.NotificationStatus
	SubmittedAt sql.NullTime
	CompletedAt sql.NullTime
	Skip        int32
	Take        int32
}
