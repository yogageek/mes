package logic

import (
	"database/sql"
	"material-management/db"
)

type Postgres struct {
	sqlDB *sql.DB
}

func NewPostgres(sqldb *sql.DB) (*Postgres, error) {
	return &Postgres{
		sqlDB: sqldb,
	}, nil
}
func newPGService() *Postgres {
	pgClient := db.CreatePGClient()
	pg, _ := NewPostgres(pgClient)
	return pg
}
