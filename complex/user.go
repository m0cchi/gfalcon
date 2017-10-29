package complex

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/m0cchi/gfalcon/model"
)

func createUser(tx *sqlx.Tx, teamIID uint32, userID string, password string) (_ *model.User, error error) {
	user, err := model.CreateUser(tx, teamIID, userID)
	if err != nil {
		return nil, err
	}

	err = user.UpdatePassword(tx, password)
	if err != nil {
		return nil, err
	}

	return user, err
}

func CreateUser(db *sqlx.DB, teamIID uint32, userID string, password string) (_ *model.User, error error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := recover(); err != nil {
			error = errors.New("failed create user")
			tx.Rollback()
		} else if error != nil {
			tx.Rollback()
		} else {
			error = tx.Commit()
		}
	}()

	return createUser(tx, teamIID, userID, password)
}
