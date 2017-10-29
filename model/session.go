package model

import (
	"time"
)

type Session struct {
	SessionID  string    `db:"session"`
	UpdateDate time.Time `db:"update_date"`
}
