package store

import (
	"context"
	"fmt"

	"github.com/aksbuzz/library-project/store/db"
	"github.com/jackc/pgx/v5"
)

type DB struct {
	*db.Queries
	*pgx.Conn
}

func NewDB(conn *pgx.Conn) *DB {
	return &DB{
		Queries: db.New(conn),
		Conn:    conn,
	}
}

func (d *DB) ExecTransaction(ctx context.Context, f func(qtx *db.Queries) error) error {
	tx, err := d.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := d.WithTx(tx)

	err = f(qtx)
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("transaction error: %w, rollback error: %w", err, rollbackErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
