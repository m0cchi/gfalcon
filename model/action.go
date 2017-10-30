package model

import (
	"github.com/go-sql-driver/mysql"
	"github.com/m0cchi/gfalcon"
)

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
