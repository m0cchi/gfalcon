package model

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/m0cchi/gfalcon"
)

const SqlGetRoleByID = "SELECT `iid`, `team_iid`, `id` FROM `roles` WHERE `team_iid` = :team_iid and `id` = :role_id"

const SqlCreateRole = "INSERT INTO `roles` (`team_iid`,`id`) VALUE (:team_iid, :role_id)"

const SqlDeleteRoleByIID = "DELETE FROM `roles` WHERE `iid` = :role_iid"
const SqlDeleteRoleByID = "DELETE FROM `roles` WHERE `team_iid` = :team_iid and `id` = :role_id"

type Role struct {
	IID     uint32 `db:"iid"`
	TeamIID uint32 `db:"team_iid"`
	ID      string `db:"id"`
}

func GetRole(db gfsql.DB, teamIID uint32, roleID string) (*Role, error) {
	stmt, err := db.PrepareNamed(SqlGetRoleByID)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	role := &Role{}
	err = stmt.Get(role, map[string]interface{}{"team_iid": teamIID, "role_id": roleID})
	return role, err
}

func CreateRole(db gfsql.DB, teamIID uint32, roleID string) (*Role, error) {
	stmt, err := db.PrepareNamed(SqlCreateRole)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	args := map[string]interface{}{"team_iid": teamIID, "role_id": roleID}
	result, err := stmt.Exec(args)
	role := &Role{IID: 0, TeamIID: teamIID, ID: roleID}

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == gfsql.ErrCodeDuplicateEntry {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	if c, err := result.LastInsertId(); err == nil {
		role.IID = uint32(c)
	}

	return role, err
}

func DeleteRoleByIID(db gfsql.DB, roleIID uint32) error {
	stmt, err := db.PrepareNamed(SqlDeleteRoleByIID)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"role_iid": roleIID}
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

func DeleteRoleByID(db gfsql.DB, teamIID uint32, roleID string) error {
	stmt, err := db.PrepareNamed(SqlDeleteRoleByID)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"team_iid": teamIID, "role_id": roleID}
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

func (role *Role) Delete(db gfsql.DB) error {
	if role == nil {
		return errors.New("not specify role")
	}

	if role.IID != 0 {
		return DeleteRoleByIID(db, role.IID)
	} else if role.ID != "" && role.TeamIID > 0 {
		return DeleteRoleByID(db, role.TeamIID, role.ID)
	}

	return errors.New("not specify role")
}
