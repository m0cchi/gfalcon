package model

import (
	"errors"
	"github.com/m0cchi/gfalcon"
)

const SqlGetActionLink = "SELECT `action_iid`, `user_iid`, `count` FROM `action_links` WHERE `action_iid` = :action_iid and `user_iid` = :user_iid"

// support role
const SqlDecrementActionLinkCount = "UPDATE `action_links` SET count = count - 1 WHERE `action_iid` = :action_iid and `user_iid` = :user_iid"

const SqlUpsertActionLink = "INSERT INTO `action_links` (`action_iid`, `user_iid`) VALUE (:action_iid, :user_iid) ON DUPLICATE KEY UPDATE count = count + 1"

const SqlCreateActionLink = "INSERT `action_links` (`action_iid`, `user_iid`) VALUE (:action_iid, :user_iid)"

const SqlDeleteActionLink = "DELETE FROM `action_links` WHERE `action_iid` = :action_iid and `user_iid` = :user_iid"

type ActionLink struct {
	ActionIID uint32 `db:"action_iid"`
	UserIID   uint32 `db:"user_iid"`
	Count     uint32 `db:"count"`
}

func upsertActionLink(db gfsql.DB, actionIID uint32, userIID uint32) error {
	stmt, err := db.PrepareNamed(SqlUpsertActionLink)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"action_iid": actionIID, "user_iid": userIID}
	_, err = stmt.Exec(args)

	return err
}

func UpsertActionLink(db gfsql.DB, action *Action, user *User) error {
	if action == nil || action.IID == 0 {
		return errors.New("not specify actionIID")
	}
	if user == nil || user.IID == 0 {
		return errors.New("not specify userIID")
	}
	return upsertActionLink(db, action.IID, user.IID)
}

func createActionLink(db gfsql.DB, actionIID uint32, userIID uint32) error {
	stmt, err := db.PrepareNamed(SqlCreateActionLink)
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
	stmt, err := db.PrepareNamed(SqlGetActionLink)
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
	stmt, err := db.PrepareNamed(SqlDecrementActionLinkCount)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"action_iid": actionIID, "user_iid": userIID}
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

func deleteActionLinkByIID(db gfsql.DB, actionIID uint32, userIID uint32) error {
	stmt, err := db.PrepareNamed(SqlDeleteActionLink)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"action_iid": actionIID, "user_iid": userIID}
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
