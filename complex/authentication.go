package complex

import (
	"errors"
	"github.com/m0cchi/gfalcon"
	"github.com/m0cchi/gfalcon/model"
	"github.com/m0cchi/gfalcon/util"
	"time"
)

const LengthOfSession = 44

const SqlGetSessionByUser = "SELECT `session`, `update_date`, `user_iid` FROM `sessions` WHERE `user_iid` = :user_iid"

const SqlUpsertSessions = "INSERT INTO `sessions` (`user_iid`,`session`) VALUES (:user_iid, :session) ON DUPLICATE KEY UPDATE `session` = :session, `update_date` = CURRENT_TIMESTAMP"

func getSessionID(db gfsql.DB, user *model.User) (*model.Session, error) {
	session := &model.Session{}
	stmt, err := db.PrepareNamed(SqlGetSessionByUser)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	args := map[string]interface{}{"user_iid": user.IID}
	err = stmt.Get(session, args)
	return session, err
}

func updateSession(db gfsql.DB, user *model.User, sessionID string) error {
	stmt, err := db.PrepareNamed(SqlUpsertSessions)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"user_iid": user.IID, "session": sessionID}
	result, err := stmt.Exec(args)
	if err != nil {
		return err
	}

	c, err := result.RowsAffected()
	if c != 1 && c != 2 {
		return errors.New("failed to update session")
	}

	return err
}

func createNewSession(user *model.User) *model.Session {
	session := &model.Session{UserIID: user.IID}
	session.SessionID = util.GenerateSessionID(LengthOfSession) // default size
	session.UpdateDate = time.Now().UTC()
	return session
}

func AuthenticateWithPassword(db gfsql.DB, user *model.User, password string) (*model.Session, error) {
	err := user.MatchPassword(db, password)
	if err != nil {
		return nil, err
	}
	// need userIID
	if user.IID < 1 {
		if user.TeamIID > 0 {
			user, err = model.GetUser(db, user.TeamIID, user.ID)
		} else {
			err = errors.New("not specify user")
		}
		if err != nil {
			return nil, err
		}
	}

	session, err := getSessionID(db, user)
	if err != nil || session == nil || session.Validate() != nil {
		session = createNewSession(user)
	}

	err = updateSession(db, user, session.SessionID)
	return session, err
}
