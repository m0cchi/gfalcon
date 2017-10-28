package model

import (
	"github.com/m0cchi/gfalcon"
)

const SQL_GET_TEAM_BY_ID = "SELECT iid, id FROM `teams` WHERE `id` = :team_id"

type Team struct {
	IID uint32 `db:"iid"`
	ID  string `db:"id"`
}

func GetTeam(db gfsql.DB, teamID string) (*Team, error) {
	stmt, err := db.PrepareNamed(SQL_GET_TEAM_BY_ID)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	team := &Team{}
	args := map[string]interface{}{"team_id": teamID}
	err = stmt.Get(team, args)
	return team, err
}
