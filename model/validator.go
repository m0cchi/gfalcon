package model

import (
	"errors"
	"time"
)

// 7 days
const EXPIRATION_INTERVAL = 7 * 24 * -1

func (session *Session) Validate() error {
	if session == nil {
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
