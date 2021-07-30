package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	balancesTable = "balances"
	freezeTable   = "freezes"
)

func NewPostgresDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", "host=localhost port=5436 user=postgres dbname=billing_db password=v&487fnd4jbvf8 sslmode=disable")

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
