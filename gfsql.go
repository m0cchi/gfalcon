package gfsql

import (
	"github.com/jmoiron/sqlx"
)

// error code of mysql
const ErrCodeDuplicateEntry = 1062

type DB interface {
	PrepareNamed(sql string) (*sqlx.NamedStmt, error)
}
