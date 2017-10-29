package complex

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/m0cchi/gfalcon/model"
)

func CreateTeam(db *sqlx.DB, teamID string, adminID string, password string) (_ *model.Team, _ *model.User, error error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err := recover(); err != nil {
			error = errors.New("failed create team")
			tx.Rollback()
		} else if error != nil {
			tx.Rollback()
		} else {
			error = tx.Commit()
		}
	}()

	team, err := model.CreateTeam(tx, teamID)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	user, err := createUser(tx, team.IID, adminID, password)

	return team, user, err
}
