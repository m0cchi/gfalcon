package complex

import (
	"database/sql"
	"errors"
	"github.com/m0cchi/gfalcon"
	"github.com/m0cchi/gfalcon/model"
	"github.com/m0cchi/gfalcon/util"
	"time"
)

// MaxChallenge is Number of times to challenge to fetch session
const MaxChallenge = 10

const LengthOfSession = 44

const SqlGetSessionByUser = "SELECT `session`, `update_date`, `user_iid` FROM `sessions` WHERE `user_iid` = :user_iid"

const SqlUpsertSessions = "INSERT INTO `sessions` (`user_iid`,`session`) VALUES (:user_iid, :session) ON DUPLICATE KEY UPDATE `session` = :session"

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
	if c != 1 {
		return errors.New("failed to delete")
	}

	return err
}

func createNewSession(db gfsql.DB) (*model.Session, error) {
	session := &model.Session{}
	i := 0
	for ; i < MaxChallenge; i++ {
		session.SessionID = util.GenerateSessionID(LengthOfSession) // default size
		s, err := model.GetSession(db, session.SessionID)
		if err == sql.ErrNoRows {
			break
		} else if err == nil && s.Validate() != nil {
			err = s.Delete(db)
			if err == nil {
				break
			}
		}
	}
	if i == MaxChallenge {
		return nil, errors.New("Oh crap, session id was exhausted")
	}

	session.UpdateDate = time.Now().UTC() // maybe unused...
	return session, nil
}

func AuthenticateWithPassword(db gfsql.DB, user *model.User, password string) (*model.Session, error) {
	err := user.MatchPassword(db, password)
	if err != nil {
		return nil, err
	}

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
		session, err = createNewSession(db)
		if err != nil {
			return nil, err
		}
	}
	err = updateSession(db, user, session.SessionID)
	return session, err
}
