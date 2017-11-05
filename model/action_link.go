package model

import (
	"errors"
	"github.com/m0cchi/gfalcon"
)

const SQL_GET_ACTION_LINK = "SELECT `action_iid`, `user_iid`, `count` FROM `action_links` WHERE `action_iid` = :action_iid and `user_iid` = :user_iid"

// support role
const SQL_INCREMENT_ACTION_LINK_COUNT = "UPDATE `action_links` SET count = count + 1 WHERE `action_iid` = :action_iid and `user_iid` = :user_iid"

// support role
const SQL_DECREMENT_ACTION_LINK_COUNT = "UPDATE `action_links` SET count = count - 1 WHERE `action_iid` = :action_iid and `user_iid` = :user_iid"

const SQL_CREATE_ACTION_LINK = "INSERT `action_links` (`action_iid`, `user_iid`) VALUE (:action_iid, :user_iid)"

const SQL_DELETE_ACTION_LINK = "DELETE FROM `action_links` WHERE `action_iid` = :action_iid and `user_iid` = :user_iid"

type ActionLink struct {
	ActionIID uint32 `db:"action_iid"`
	UserIID   uint32 `db:"user_iid"`
	Count     uint32 `db:"count"`
}

/*
support role
func incrementActionLinkCount(db gfsql.DB, actionIID uint32, userIID uint32) error {
	stmt, err := db.PrepareNamed(SQL_INCREMENT_ACTION_LINK_COUNT)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"action_iid": actionIID, "user_iid": userIID}
	_, err = stmt.Exec(args)
	return err
}
*/

func createActionLink(db gfsql.DB, actionIID uint32, userIID uint32) error {
	stmt, err := db.PrepareNamed(SQL_CREATE_ACTION_LINK)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"action_iid": actionIID, "user_iid": userIID}
	_, err = stmt.Exec(args)

	return err
}

func CreateActionLink(db gfsql.DB, action *Action, user *User) error {
	if action == nil || action.IID == 0 {
		return errors.New("not specify actionIID")
	}
	if user == nil || user.IID == 0 {
		return errors.New("not specify userIID")
	}
	return createActionLink(db, action.IID, user.IID)
}

func getActionLink(db gfsql.DB, actionIID uint32, userIID uint32) (*ActionLink, error) {
	stmt, err := db.PrepareNamed(SQL_GET_ACTION_LINK)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	args := map[string]interface{}{"action_iid": actionIID, "user_iid": userIID}
	actionLink := &ActionLink{}
	err = stmt.Get(actionLink, args)
	return actionLink, err
}

func GetActionLink(db gfsql.DB, action *Action, user *User) (*ActionLink, error) {
	if action == nil || action.IID == 0 {
		return nil, errors.New("not specify actionIID")
	}
	if user == nil || user.IID == 0 {
		return nil, errors.New("not specify userIID")
	}
	return getActionLink(db, action.IID, user.IID)
}

func decrementActionLinkCount(db gfsql.DB, actionIID uint32, userIID uint32) error {
	stmt, err := db.PrepareNamed(SQL_DECREMENT_ACTION_LINK_COUNT)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"action_iid": actionIID, "user_iid": userIID}
	_, err = stmt.Exec(args)
	return err
}

func deleteActionLinkByIID(db gfsql.DB, actionIID uint32, userIID uint32) error {
	stmt, err := db.PrepareNamed(SQL_DELETE_ACTION_LINK)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"action_iid": actionIID, "user_iid": userIID}
	_, err = stmt.Exec(args)
	return err
}

func deleteActionLink(db gfsql.DB, actionIID uint32, userIID uint32) error {
	actionLink, err := getActionLink(db, actionIID, userIID)
	if err != nil {
		return errors.New("not found action link")
	}
	if actionLink.Count == 1 {
		// delete
		return deleteActionLinkByIID(db, actionIID, userIID)
	}
	//decrement
	return decrementActionLinkCount(db, actionIID, userIID)
}

func DeleteActionLink(db gfsql.DB, action *Action, user *User) error {
	if action == nil || action.IID == 0 {
		return errors.New("not specify actionIID")
	}
	if user == nil || user.IID == 0 {
		return errors.New("not specify userIID")
	}
	return deleteActionLink(db, action.IID, user.IID)
}
