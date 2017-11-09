package model

import (
	"github.com/m0cchi/gfalcon"
	"time"
)

const SqlGetSession = "SELECT `session`, `update_date`, `user_iid` FROM `sessions` WHERE `session` = :session_id"

type Session struct {
	SessionID  string    `db:"session"`
	UpdateDate time.Time `db:"update_date"`
	UserIID    uint32    `db:"user_iid"`
}

func GetSession(db gfsql.DB, sessionID string) (*Session, error) {
	session := &Session{}
	stmt, err := db.PrepareNamed(SqlGetSession)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	args := map[string]interface{}{"session_id": sessionID}
	err = stmt.Get(session, args)
	return session, err
}
