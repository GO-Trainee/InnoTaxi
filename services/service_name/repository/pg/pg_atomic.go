package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type PgAtomicRepository interface {
	Start(context.Context) (*sqlx.Tx, error)
	Finish(*sqlx.Tx) error
	Abort(*sqlx.Tx) error
}
type pgAtomicRepository struct {
	db *sqlx.DB
}

func NewAtomic(db *sqlx.DB) PgAtomicRepository {
	return &pgAtomicRepository{db: db}
}

func (r *pgAtomicRepository) Start(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *pgAtomicRepository) Finish(tx *sqlx.Tx) error {
	if tx == nil {
		return errors.New("transaction is not started")
	}
	return tx.Commit()
}

func (r *pgAtomicRepository) Abort(tx *sqlx.Tx) error {
	if tx == nil {
		return errors.New("transaction is not started")
	}
	err := tx.Rollback()
	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		return err
	}
	return nil
}
