package model

import (
	"errors"
	"time"
)

// 7 days
const ExpirationInterval = 7 * 24 * -1

func (session *Session) Validate() error {
	if session == nil {
		return errors.New("session is nil")
	}

	if session.SessionID == "" {
		return errors.New("SessionID is required field")
	}

	sub := time.Since(session.UpdateDate)
	if sub.Hours() <= ExpirationInterval {
		return errors.New("session has expired")
	}
	return nil
}

func (actionLink *ActionLink) Validate() error {
	if actionLink == nil {
		return errors.New("action link is nil")
	}

	if actionLink.Count < 1 {
		return errors.New("invalid count")
	}

	if actionLink.ActionIID < 0 {
		return errors.New("invalid actionIID")
	}

	if actionLink.UserIID < 0 {
		return errors.New("invalid userIID")
	}
	return nil
}
