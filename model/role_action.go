package model

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/m0cchi/gfalcon"
)

const SqlGetRoleActions = "SELECT `role_iid`, `action_iid` FROM `role_actions` WHERE `role_iid` = :role_iid"

const SqlCreateRoleAction = "INSERT INTO `role_actions` (`role_iid`, `action_iid`) VALUE (:role_iid, :action_iid)"

const SqlDeleteRoleAction = "DELETE FROM `role_actions` WHERE `action_iid` = :action_iid and `role_iid` = :role_iid"

type RoleAction struct {
	RoleIID   uint32 `db:"role_iid"`
	ActionIID uint32 `db:"action_iid"`
}

func getRoleActions(db gfsql.DB, roleIID uint32) ([]RoleAction, error) {
	stmt, err := db.PrepareNamed(SqlGetRoleActions)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	args := map[string]interface{}{"role_iid": roleIID}
	roleActions := make([]RoleAction, 0, 10)
	err = stmt.Select(&roleActions, args)

	return roleActions, err
}

func GetRoleActions(db gfsql.DB, role *Role) ([]RoleAction, error) {
	if role == nil || role.IID == 0 {
		return nil, errors.New("not specify roleIID")
	}
	return getRoleActions(db, role.IID)
}

func createRoleAction(db gfsql.DB, roleIID uint32, actionIID uint32) error {
	stmt, err := db.PrepareNamed(SqlCreateRoleAction)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"action_iid": actionIID, "role_iid": roleIID}
	_, err = stmt.Exec(args)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == gfsql.ErrCodeDuplicateEntry {
				return ErrDuplicate
			}
		}
		return err
	}
	return err
}

func CreateRoleAction(db gfsql.DB, role *Role, action *Action) error {
	if role == nil || role.IID == 0 {
		return errors.New("not specify roleIID")
	}
	if action == nil || action.IID == 0 {
		return errors.New("not specify actionIID")
	}

	return createRoleAction(db, role.IID, action.IID)
}

func deleteRoleAction(db gfsql.DB, roleIID uint32, actionIID uint32) error {
	stmt, err := db.PrepareNamed(SqlDeleteRoleAction)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"action_iid": actionIID, "role_iid": roleIID}
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

func DeleteRoleAction(db gfsql.DB, role *Role, action *Action) error {
	if action == nil || action.IID == 0 {
		return errors.New("not specify actionIID")
	}
	if role == nil || role.IID == 0 {
		return errors.New("not specify roleIID")
	}
	return deleteRoleAction(db, role.IID, action.IID)
}
