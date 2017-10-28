package complex

import (
	"errors"
	"time"
)

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
