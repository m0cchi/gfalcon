package complex

import (
	"errors"
	"github.com/jmoiron/sqlx"

	"github.com/m0cchi/gfalcon/model"
)

func CreateUser(db *sqlx.DB, teamIID uint32, userID string, password string) (_ *model.User, error error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := recover(); err != nil {
			error = errors.New("failed create user")
			tx.Rollback()
		}
	}()
	user, err := model.CreateUser(tx, teamIID, userID)
	if err != nil {
		tx.Rollback()
		return user, err
	}

	err = user.UpdatePassword(tx, password)
	if err != nil {
		tx.Rollback()
		return user, err
	}

	tx.Commit()
	return user, err
}
