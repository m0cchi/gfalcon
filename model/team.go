package model

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/m0cchi/gfalcon"
)

const SQL_GET_TEAM_BY_ID = "SELECT iid, id FROM `teams` WHERE `id` = :team_id"
const SQL_CREATE_TEAM = "INSERT INTO `teams` (`id`) VALUE (:team_id)"

const SQL_DELETE_TEAM_BY_IID = "DELETE FROM `teams` WHERE `iid` = :team_iid"
const SQL_DELETE_TEAM_BY_ID = "DELETE FROM `teams` WHERE `id` = :team_id"

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
			if mysqlErr.Number == gfsql.ErrCodeDuplicateEntry {
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

func DeleteTeamByIID(db gfsql.DB, teamIID uint32) error {
	stmt, err := db.PrepareNamed(SQL_DELETE_TEAM_BY_IID)
	if err != nil {
		return err
	}
	defer stmt.Close()

	args := map[string]interface{}{"team_iid": teamIID}
	_, err = stmt.Exec(args)
	return err
}

func DeleteTeamByID(db gfsql.DB, teamID string) error {
	stmt, err := db.PrepareNamed(SQL_DELETE_TEAM_BY_ID)
	if err != nil {
		return err
	}
	defer stmt.Close()

	args := map[string]interface{}{"team_id": teamID}
	_, err = stmt.Exec(args)
	return err
}

func (team *Team) Delete(db gfsql.DB) error {
	if team == nil {
		return errors.New("not specify team")
	}

	if team.IID != 0 {
		return DeleteTeamByIID(db, team.IID)
	} else if team.ID != "" {
		return DeleteTeamByID(db, team.ID)
	}

	return errors.New("not specify team")
}
