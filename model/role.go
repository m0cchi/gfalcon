package model

import (
	"github.com/m0cchi/gfalcon"
)

const SqlGetRoleByID = "SELECT `iid`, `team_iid`, `id` FROM `roles` WHERE `team_iid` = :team_iid and `id` = :role_id"

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
