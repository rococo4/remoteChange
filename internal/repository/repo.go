package repository

import "github.com/jmoiron/sqlx"

type Repo struct {
	Db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{Db: db}
}
