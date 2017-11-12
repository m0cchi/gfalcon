package model

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/m0cchi/gfalcon"
)

const SqlGetActionByID = "SELECT `iid`, `service_iid`, `id` FROM `actions` WHERE `id` = :action_id and `service_iid` = :service_iid"

const SqlCreateAction = "INSERT INTO `actions` (`service_iid`, `id`) VALUE (:service_iid, :action_id)"

const SqlDeleteActionByIID = "DELETE FROM `actions` WHERE `iid` = :action_iid and `service_iid` = :service_iid"
const SqlDeleteActionByID = "DELETE FROM `actions` WHERE `id` = :action_id and `service_iid` = :service_iid"

type Action struct {
	IID        uint32 `db:"iid"`
	ServiceIID uint32 `db:"service_iid"`
	ID         string `db:"id"`
}

func CreateAction(db gfsql.DB, serviceIID uint32, actionID string) (*Action, error) {
	stmt, err := db.PrepareNamed(SqlCreateAction)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	args := map[string]interface{}{"service_iid": serviceIID, "action_id": actionID}
	result, err := stmt.Exec(args)
	action := &Action{0, serviceIID, actionID}

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == gfsql.ErrCodeDuplicateEntry {
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
	stmt, err := db.PrepareNamed(SqlGetActionByID)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	action := &Action{}
	args := map[string]interface{}{"service_iid": serviceIID, "action_id": actionID}
	err = stmt.Get(action, args)
	return action, err
}

func DeleteActionByIID(db gfsql.DB, serviceIID uint32, actionIID uint32) error {
	stmt, err := db.PrepareNamed(SqlDeleteActionByIID)
	if err != nil {
		return err
	}
	defer stmt.Close()

	args := map[string]interface{}{"service_iid": serviceIID, "action_iid": actionIID}
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

func DeleteActionByID(db gfsql.DB, serviceIID uint32, actionID string) error {
	stmt, err := db.PrepareNamed(SqlDeleteActionByID)
	if err != nil {
		return err
	}
	defer stmt.Close()

	args := map[string]interface{}{"service_iid": serviceIID, "action_id": actionID}
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

func (action *Action) Delete(db gfsql.DB) error {
	if action == nil || action.ServiceIID < 1 {
		return errors.New("not specify action")
	}

	if action.IID != 0 {
		return DeleteActionByIID(db, action.ServiceIID, action.IID)
	} else if action.ID != "" {
		return DeleteActionByID(db, action.ServiceIID, action.ID)
	}

	return errors.New("not specify action")
}
