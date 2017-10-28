package complex

import (
	"errors"

	"github.com/m0cchi/gfalcon"
	"github.com/m0cchi/gfalcon/model"
	"github.com/m0cchi/gfalcon/util"
	"time"
)

const LENGTH_OF_SESSION = 44

// 7 days
const EXPIRATION_INTERVAL = 1 * 24 * -1

const SQL_GET_SESSION_BY_USER = "SELECT `session`, `update_date` FROM `sessions` WHERE `team_iid` = :team_iid and `user_iid` = :user_iid"

const SQL_UPSERT_SESSIONS = "INSERT INTO `sessions` (`team_iid`,`user_iid`,`session`) VALUES (:team_iid, :user_iid, :session) ON DUPLICATE KEY UPDATE `session` = :session"

type Session struct {
	SessionID  string    `db:"session"`
	UpdateDate time.Time `db:"update_date"`
}

func getSessionID(db gfsql.DB, user *model.User) (*Session, error) {
	session := &Session{}
	stmt, err := db.PrepareNamed(SQL_GET_SESSION_BY_USER)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	args := map[string]interface{}{"team_iid": user.TeamIID, "user_iid": user.IID}
	err = stmt.Get(session, args)
	return session, err
}

func updateSession(db gfsql.DB, user *model.User, sessionID string) error {
	stmt, err := db.PrepareNamed(SQL_UPSERT_SESSIONS)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"team_iid": user.TeamIID, "user_iid": user.IID, "session": sessionID}
	_, err = stmt.Exec(args)

	return err
}

func (session *Session) Validate() error {
	if session != nil {
		return errors.New("session is nil")
	}

	if session.SessionID == "" {
		return errors.New("SessionID is required field")
	}

	sub := time.Since(session.UpdateDate)
	if sub.Hours() <= EXPIRATION_INTERVAL {
		return errors.New("session has expired")
	}
	return nil
}

func AuthenticateWithPassword(db gfsql.DB, user *model.User, password string) (*Session, error) {
	err := user.MatchPassword(db, password)
	if err != nil {
		return nil, err
	}

	session, err := getSessionID(db, user)
	if err != nil || session == nil || session.Validate() != nil {
		// new session
		session.SessionID = util.GenerateSessionID(LENGTH_OF_SESSION) // default size
		session.UpdateDate = time.Now().UTC()                         // maybe unused...
	}
	err = updateSession(db, user, session.SessionID)
	return session, err
}
