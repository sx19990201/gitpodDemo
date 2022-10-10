package sqllite

import (
	"database/sql"
)

type homeRepository struct {
	db *sql.DB
}

func NewHomeRepository(Conn *sql.DB) *homeRepository {
	return &homeRepository{
		db: Conn,
	}
}
