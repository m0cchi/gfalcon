package gfsql

import (
	"github.com/jmoiron/sqlx"
)

// error code of mysql
const ERR_CODE_DUPLICATE_ENTRY = 1062

type DB interface {
	PrepareNamed(sql string) (*sqlx.NamedStmt, error)
}
