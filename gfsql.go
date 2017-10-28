package gfsql

import (
	"github.com/jmoiron/sqlx"
)

type DB interface {
	PrepareNamed(sql string) (*sqlx.NamedStmt, error)
}
