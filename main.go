package main

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Close()
}

func recordStats(db PgxIface, userID, productID int) (err error) {
	var tx pgx.Tx
	if tx, err = db.Begin(context.Background()); err != nil {
		return err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit(context.Background())
		default:
			_ = tx.Rollback(context.Background())
		}
	}()
	sql := "UPDATE products SET views = views + 1"
	if _, err = tx.Exec(context.Background(), sql); err != nil {
		return err
	}
	sql = "INSERT INTO product_viewers (user_id, product_id) VALUES ($1, $2)"
	if _, err = tx.Exec(context.Background(), sql, userID, productID); err != nil {
		return
	}
	sql = "SELECT user_id FROM product_viewers WHERE user_id = $1 AND product_id = $2"
	result, err := tx.Query(context.Background(), sql, userID, productID)
	if err != nil {
		return
	}
	defer result.Close()
	for result.Next() {
		var id int
		if err = result.Scan(&id); err != nil {
			return
		}
		print(id)

	}

	return
}

func main() {
	// @NOTE: the real connection is not required for tests
	db, err := pgxpool.New(context.Background(), "postgres://aditya@localhost/postgres")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = recordStats(db, 1 /*some user id*/, 5 /*some product id*/); err != nil {
		panic(err)
	}
}
