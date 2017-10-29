package model

import (
	"github.com/go-sql-driver/mysql"
	"github.com/m0cchi/gfalcon"
)

const SQL_GET_TEAM_BY_ID = "SELECT iid, id FROM `teams` WHERE `id` = :team_id"
const SQL_CREATE_TEAM = "INSERT INTO `teams` (`id`) VALUE (:team_id)"

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

func CreateTeam(db gfsql.DB, teamID string) (*Team, error) {
	stmt, err := db.PrepareNamed(SQL_CREATE_TEAM)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	args := map[string]interface{}{"team_id": teamID}
	result, err := stmt.Exec(args)
	team := &Team{
		IID: 0,
		ID:  teamID}

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == gfsql.ERR_CODE_DUPLICATE_ENTRY {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	if c, err := result.LastInsertId(); err == nil {
		team.IID = uint32(c)
	}

	return team, err
}
