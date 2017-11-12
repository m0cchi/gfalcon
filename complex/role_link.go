package complex

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/m0cchi/gfalcon/model"
)

func createRoleLink(tx *sqlx.Tx, role *model.Role, user *model.User) error {
	err := model.CreateRoleLink(tx, role, user)
	if err != nil {
		return err
	}

	roleActions, err := model.GetRoleActions(tx, role)
	if err != nil {
		return err
	}
	for _, roleAction := range roleActions {
		action := &model.Action{IID: roleAction.ActionIID}
		err := model.CreateActionLink(tx, action, user)
		if err != nil {
			return err
		}
	}
	return err
}

func CreateRoleLink(db sqlx.DB, role *model.Role, user *model.User) (error error) {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			error = errors.New("failed create role link")
			tx.Rollback()
		} else if error != nil {
			tx.Rollback()
		} else {
			error = tx.Commit()
		}
	}()
	return createRoleLink(tx, role, user)
}
