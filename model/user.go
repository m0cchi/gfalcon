package model

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/m0cchi/gfalcon"
	"golang.org/x/crypto/bcrypt"
)

const SQL_GET_USER_BY_ID = "SELECT `iid`, `team_iid`, `id` FROM `users` WHERE `team_iid` = :team_iid and `id` = :user_id"
const SQL_UPSERT_PASSWORD_INSERT_BASE = "INSERT INTO `passwords` (`user_iid`,`password`) SELECT `users`.`iid`, :password as password FROM `users` WHERE "
const SQL_UPSERT_PASSWORD_UPDATE_BASE = "ON DUPLICATE KEY UPDATE `passwords`.`password` = :password"
const SQL_UPSERT_PASSWORD_BY_IID = SQL_UPSERT_PASSWORD_INSERT_BASE + "`users`.`iid` = :user_iid " + SQL_UPSERT_PASSWORD_UPDATE_BASE
const SQL_UPSERT_PASSWORD_BY_ID = SQL_UPSERT_PASSWORD_INSERT_BASE + "`users`.`team_iid` = :team_iid and `users`.`id` = :user_id " + SQL_UPSERT_PASSWORD_UPDATE_BASE

const SQL_GET_PASSWORD_BY_IID = "SELECT `password` FROM `passwords` WHERE `user_iid` = :user_iid"
const SQL_GET_PASSWORD_BY_ID = "SELECT `password` FROM `passwords`, (SELECT `user_iid` FROM `users` WHERE `team_iid` = :team_iid and `id` = :user_id) as `filtered` WHERE `passwords`.`user_iid` = `filtered`.`user_iid`"

const SQL_CREATE_USER = "INSERT INTO `users` (`team_iid`,`id`) VALUE (:team_iid, :user_id)"

const SQL_DELETE_USER_BY_IID = "DELETE FROM `users` WHERE `iid` = :user_iid"
const SQL_DELETE_USER_BY_ID = "DELETE FROM `users` WHERE `team_iid` = :team_iid and `id` = :user_id"

type User struct {
	IID     uint32 `db:"iid"`
	TeamIID uint32 `db:"team_iid"`
	ID      string `db:"id"`
}

type Password struct {
	Password string `db:"password"`
}

func toHash(password string) string {
	converted, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(converted)
}

func (user *User) MatchPassword(db gfsql.DB, password string) error {
	var stmt *sqlx.NamedStmt
	var err error
	var args interface{}

	if user.IID != 0 {
		stmt, err = db.PrepareNamed(SQL_GET_PASSWORD_BY_IID)
		args = map[string]interface{}{"user_iid": user.IID}
	} else if user.ID != "" && user.TeamIID > 0 {
		stmt, err = db.PrepareNamed(SQL_GET_PASSWORD_BY_ID)
		args = map[string]interface{}{"team_iid": user.TeamIID, "user_id": user.ID}
	} else {
		return errors.New("not specify IID or ID")
	}
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedpassword := Password{}
	err = stmt.Get(&hashedpassword, args)

	if err != nil || hashedpassword.Password == "" {
		// missing user's pasword
		return errors.New("Unmatch password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedpassword.Password), []byte(password))
	if err != nil {
		return errors.New("Unmatch password")
	}
	return nil

}

func (user *User) UpdatePassword(db gfsql.DB, password string) error {
	var stmt *sqlx.NamedStmt
	var err error
	var args interface{}

	if user == nil {
		return errors.New("not specify IID or ID")
	}

	if user.IID != 0 {
		stmt, err = db.PrepareNamed(SQL_UPSERT_PASSWORD_BY_IID)
		args = map[string]interface{}{"user_iid": user.IID, "password": toHash(password)}
	} else if user.ID != "" && user.TeamIID > 0 {
		stmt, err = db.PrepareNamed(SQL_UPSERT_PASSWORD_BY_ID)
		args = map[string]interface{}{"team_iid": user.TeamIID, "user_id": user.ID, "password": toHash(password)}
	} else {
		return errors.New("not specify IID or ID")
	}
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(args)
	return err
}

func GetUser(db gfsql.DB, teamIID uint32, userID string) (*User, error) {
	stmt, err := db.PrepareNamed(SQL_GET_USER_BY_ID)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	user := &User{}
	err = stmt.Get(user, map[string]interface{}{"team_iid": teamIID, "user_id": userID})
	return user, err
}

func CreateUser(db gfsql.DB, teamIID uint32, userID string) (*User, error) {
	stmt, err := db.PrepareNamed(SQL_CREATE_USER)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	args := map[string]interface{}{"team_iid": teamIID, "user_id": userID}
	result, err := stmt.Exec(args)
	user := &User{0, teamIID, userID}

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == gfsql.ErrCodeDuplicateEntry {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	if c, err := result.LastInsertId(); err == nil {
		user.IID = uint32(c)
	}

	return user, err
}

func DeleteUserByIID(db gfsql.DB, userIID uint32) error {
	stmt, err := db.PrepareNamed(SQL_DELETE_USER_BY_IID)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"user_iid": userIID}
	_, err = stmt.Exec(args)

	return err
}

func DeleteUserByID(db gfsql.DB, teamIID uint32, userID string) error {
	stmt, err := db.PrepareNamed(SQL_DELETE_USER_BY_ID)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"team_iid": teamIID, "user_id": userID}
	_, err = stmt.Exec(args)

	return err
}

func (user *User) Delete(db gfsql.DB) error {
	if user == nil {
		return errors.New("not specify user")
	}

	if user.IID != 0 {
		return DeleteUserByIID(db, user.IID)
	} else if user.ID != "" && user.TeamIID > 0 {
		return DeleteUserByID(db, user.TeamIID, user.ID)
	}

	return errors.New("not specify user")
}
