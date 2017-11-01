package model

import (
	"github.com/go-sql-driver/mysql"
	"github.com/m0cchi/gfalcon"
)

const SQL_GET_ACTION_BY_ID = "SELECT `iid`, `service_iid`, `id` FROM `actions` WHERE `id` = :action_id and `service_iid` = :service_iid"

const SQL_CREATE_ACTION = "INSERT INTO `actions` (`service_iid`, `id`) VALUE (:service_iid, :action_id)"

type Action struct {
	IID        uint32 `db:"iid"`
	ServiceIID uint32 `db:"service_iid"`
	ID         string `db:"id"`
}

func CreateAction(db gfsql.DB, serviceIID uint32, actionID string) (*Action, error) {
	stmt, err := db.PrepareNamed(SQL_CREATE_ACTION)
	if err != nil {
		return nil, err
	}

	args := map[string]interface{}{"service_iid": serviceIID, "action_id": actionID}
	result, err := stmt.Exec(args)
	action := &Action{0, serviceIID, actionID}

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == gfsql.ERR_CODE_DUPLICATE_ENTRY {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	if c, err := result.LastInsertId(); err == nil {
		action.IID = uint32(c)
	}

	return action, err
}

func GetAction(db gfsql.DB, serviceIID uint32, actionID string) (*Action, error) {
	stmt, err := db.PrepareNamed(SQL_GET_ACTION_BY_ID)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	action := &Action{}
	args := map[string]interface{}{"service_iid": serviceIID, "action_id": actionID}
	err = stmt.Get(action, args)
	return action, err
}

