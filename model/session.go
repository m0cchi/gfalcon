package model

import (
	"errors"
	"github.com/m0cchi/gfalcon"
	"time"
)

const SqlGetSession = "SELECT `session`, `update_date`, `user_iid` FROM `sessions` WHERE `session` = :session_id and `user_iid` = :user_iid"
const SqlDeleteSession = "DELETE FROM `sessions` WHERE `session` = :session_id and `user_iid` = :user_iid"

type Session struct {
	SessionID  string    `db:"session"`
	UpdateDate time.Time `db:"update_date"`
	UserIID    uint32    `db:"user_iid"`
}

func GetSession(db gfsql.DB, userIID uint32, sessionID string) (*Session, error) {
	session := &Session{}
	stmt, err := db.PrepareNamed(SqlGetSession)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	args := map[string]interface{}{"session_id": sessionID, "user_iid": userIID}
	err = stmt.Get(session, args)
	return session, err
}

func deleteSession(db gfsql.DB, userIID uint32, sessionID string) error {
	stmt, err := db.PrepareNamed(SqlDeleteSession)
	if err != nil {
		return err
	}
	defer stmt.Close()

	args := map[string]interface{}{"session_id": sessionID, "user_iid": userIID}
	result, err := stmt.Exec(args)
	if err != nil {
		return err
	}

	c, err := result.RowsAffected()
	if c != 1 {
		return errors.New("failed to delete")
	}

	return err
}

func (session *Session) Delete(db gfsql.DB) error {
	if session == nil {
		return errors.New("not specify session")
	}
	if session.SessionID == "" {
		return errors.New("not specify SessionID")
	}
	if session.UserIID <= 0 {
		return errors.New("not specify UserIIDD")
	}
	return deleteSession(db, session.UserIID, session.SessionID)
}
