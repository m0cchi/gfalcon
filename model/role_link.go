package model

import (
	"errors"
	"github.com/m0cchi/gfalcon"
)

type RoleLink struct {
	RoleIID uint32 `db:"role_iid"`
	UserIID uint32 `db:"user_iid"`
}

const SqlCreateRoleLink = "INSERT `role_links` (`role_iid`, `user_iid`) VALUE (:role_iid, :user_iid)"

const SqlDeleteRoleLink = "DELETE FROM `role_links` WHERE `role_iid` = :role_iid and `user_iid` = :user_iid"

func createRoleLink(db gfsql.DB, roleIID uint32, userIID uint32) error {
	stmt, err := db.PrepareNamed(SqlCreateRoleLink)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"role_iid": roleIID, "user_iid": userIID}
	_, err = stmt.Exec(args)

	return err
}

func CreateRoleLink(db gfsql.DB, role *Role, user *User) error {
	if role == nil || role.IID == 0 {
		return errors.New("not specify roleIID")
	}
	if user == nil || user.IID == 0 {
		return errors.New("not specify userIID")
	}
	return createRoleLink(db, role.IID, user.IID)
}

func deleteRoleLinkByIID(db gfsql.DB, roleIID uint32, userIID uint32) error {
	stmt, err := db.PrepareNamed(SqlDeleteRoleLink)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"role_iid": roleIID, "user_iid": userIID}
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

func DeleteRoleLink(db gfsql.DB, role *Role, user *User) error {
	if role == nil || role.IID == 0 {
		return errors.New("not specify roleIID")
	}
	if user == nil || user.IID == 0 {
		return errors.New("not specify userIID")
	}
	return deleteRoleLinkByIID(db, role.IID, user.IID)
}
