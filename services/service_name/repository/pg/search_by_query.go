func (r *Repository) SearchNotification(ctx context.Context, q *repository.NotificationSearchQuery) ([]*entity.Notification, error) {
	var (
		notifications []*entity.Notification
		conditions    []string
		params        []interface{}
	)

	query := `SELECT * FROM notifications`

	if q.JobID.Valid {
		conditions = append(conditions, "job_id = ?")
		params = append(params, q.JobID.String)
	}

	if q.ListID.Valid && q.ListID.Int64 > 0 {
		conditions = append(conditions, "list_id = ?")
		params = append(params, q.ListID.Int64)
	}

	if q.Status != entity.NotificationStatusUnspecified {
		conditions = append(conditions, "status = ?")
		params = append(params, q.Status)
	}

	if q.SubmittedAt.Valid {
		conditions = append(conditions, "submitted_at >= ?")
		params = append(params, q.SubmittedAt.Time)
	}

	if q.CompletedAt.Valid {
		conditions = append(conditions, "completed_at >= ?")
		params = append(params, q.CompletedAt.Time)
	}

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	query = fmt.Sprintf("%s LIMIT ? OFFSET ?", query)
	params = append(params, q.Take, q.Skip)

	query, args, err := sqlx.In(query, params...)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	if err := r.db.SelectContext(ctx, &notifications, query, args...); err != nil {
		return nil, err
	}

	return notifications, nil
}