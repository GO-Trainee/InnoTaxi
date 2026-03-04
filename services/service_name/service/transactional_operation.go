package service

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func (s *service) transactionalOperation(ctx context.Context, txFunc func(tx *sqlx.Tx) error) error {
	tx, err := s.repo.Start(ctx)
	if err != nil {
		s.logger.WithMethod("transactionalOperation").Error(err, "msg", "failed to start transaction")
		return err
	}

	defer func() {
		if err != nil {
			if abortErr := s.repo.Abort(tx); abortErr != nil {
				s.logger.WithMethod("transactionalOperation").Error(abortErr, "msg", "failed to abort transaction")
			}
		}
	}()

	if err = txFunc(tx); err != nil {
		s.logger.WithMethod("transactionalOperation").Error(err, "msg", "failed to execute transactional operation")
		return err
	}

	if err = s.repo.Finish(tx); err != nil {
		s.logger.WithMethod("transactionalOperation").Error(err, "msg", "failed to finish transaction")
		return err
	}

	return nil
}
