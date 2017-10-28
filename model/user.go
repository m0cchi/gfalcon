package model

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/m0cchi/gfalcon"
	"golang.org/x/crypto/bcrypt"
)

const SQL_GET_USER_BY_ID = "SELECT `iid`, `team_iid`, `id` FROM `users` WHERE `team_iid` = :team_iid and `id` = :user_id"
const SQL_UPSERT_PASSWORD_INSERT_BASE = "INSERT INTO `passwords` (`team_iid`,`user_iid`,`password`) SELECT `users`.`team_iid`, `users`.`iid`, :password as password FROM `users` WHERE `users`.`team_iid` = :team_iid and "
const SQL_UPSERT_PASSWORD_UPDATE_BASE = "ON DUPLICATE KEY UPDATE `passwords`.`password` = :password"
const SQL_UPSERT_PASSWORD_BY_IID = SQL_UPSERT_PASSWORD_INSERT_BASE + "`users`.`iid` = :user_iid " + SQL_UPSERT_PASSWORD_UPDATE_BASE
const SQL_UPSERT_PASSWORD_BY_ID = SQL_UPSERT_PASSWORD_INSERT_BASE + "`users`.`id` = :user_id " + SQL_UPSERT_PASSWORD_UPDATE_BASE

const SQL_GET_PASSWORD_BY_IID = "SELECT `password` FROM `passwords` WHERE `team_iid` = :team_iid and `user_iid` = :user_iid"
const SQL_GET_PASSWORD_BY_ID = "SELECT `password` FROM `passwords`, (SELECT `user_iid` FROM `users` WHERE `team_iid` = :team_iid and `id` = :user_id) as `filtered` WHERE `passwords`.`team_iid` = :team_iid and `passwords`.`user_iid` = `filtered`.`user_iid`"

const SQL_CREATE_USER = "INSERT INTO `users` (`team_iid`,`id`) VALUE (:team_iid, :user_id)"

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

	if user.TeamIID == 0 {
		return errors.New("not specify TeamIID")
	}

	if user.IID != 0 {
		stmt, err = db.PrepareNamed(SQL_GET_PASSWORD_BY_IID)
		args = map[string]interface{}{"team_iid": user.TeamIID, "user_iid": user.IID}
	} else if user.ID != "" {
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

	if user.TeamIID == 0 {
		return errors.New("not specify TeamIID")
	}

	if user.IID != 0 {
		stmt, err = db.PrepareNamed(SQL_UPSERT_PASSWORD_BY_IID)
		args = map[string]interface{}{"team_iid": user.TeamIID, "user_iid": user.IID, "password": toHash(password)}
	} else if user.ID != "" {
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
	// TODO: make error of duplicate entry
	stmt, err := db.PrepareNamed(SQL_CREATE_USER)
	if err != nil {
		return nil, err
	}

	args := map[string]interface{}{"team_iid": teamIID, "user_id": userID}
	result, err := stmt.Exec(args)
	user := &User{0, teamIID, userID}
	if c, err := result.LastInsertId(); err == nil {
		user.IID = uint32(c)
	}

	return user, err
}

func CreateUserWithPassword(db *sqlx.DB, teamIID uint32, userID string, password string) (_ *User, error error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := recover(); err != nil {
			error = errors.New("failed create user")
			tx.Rollback()
		}
	}()
	user, err := CreateUser(tx, teamIID, userID)
	if err != nil {
		tx.Rollback()
		return user, err
	}

	err = user.UpdatePassword(tx, password)
	if err != nil {
		tx.Rollback()
		return user, err
	}

	tx.Commit()
	return user, err
}
