package model

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/m0cchi/gfalcon"
	"golang.org/x/crypto/bcrypt"
)

const SqlGetUserByID = "SELECT `iid`, `team_iid`, `id` FROM `users` WHERE `team_iid` = :team_iid and `id` = :user_id"
const SqlUpsertPasswordInsertBase = "INSERT INTO `passwords` (`user_iid`,`password`) SELECT `users`.`iid`, :password as password FROM `users` WHERE "
const SqlUpsertPasswordUpdateBase = "ON DUPLICATE KEY UPDATE `passwords`.`password` = :password"
const SqlUpsertPasswordByIID = SqlUpsertPasswordInsertBase + "`users`.`iid` = :user_iid " + SqlUpsertPasswordUpdateBase
const SqlUpsertPasswordByID = SqlUpsertPasswordInsertBase + "`users`.`team_iid` = :team_iid and `users`.`id` = :user_id " + SqlUpsertPasswordUpdateBase

const SqlGetPasswordByIID = "SELECT `password` FROM `passwords` WHERE `user_iid` = :user_iid"
const SqlGetPasswordByID = "SELECT `password` FROM `passwords`, (SELECT `user_iid` FROM `users` WHERE `team_iid` = :team_iid and `id` = :user_id) as `filtered` WHERE `passwords`.`user_iid` = `filtered`.`user_iid`"

const SqlCreateUser = "INSERT INTO `users` (`team_iid`,`id`) VALUE (:team_iid, :user_id)"

const SqlDeleteUserByIID = "DELETE FROM `users` WHERE `iid` = :user_iid"
const SqlDeleteUserByID = "DELETE FROM `users` WHERE `team_iid` = :team_iid and `id` = :user_id"

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

	if user.IID > 0 {
		stmt, err = db.PrepareNamed(SqlGetPasswordByIID)
		args = map[string]interface{}{"user_iid": user.IID}
	} else if user.ID != "" && user.TeamIID > 0 {
		stmt, err = db.PrepareNamed(SqlGetPasswordByID)
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

	if user.IID > 0 {
		stmt, err = db.PrepareNamed(SqlUpsertPasswordByIID)
		args = map[string]interface{}{"user_iid": user.IID, "password": toHash(password)}
	} else if user.ID != "" && user.TeamIID > 0 {
		stmt, err = db.PrepareNamed(SqlUpsertPasswordByID)
		args = map[string]interface{}{"team_iid": user.TeamIID, "user_id": user.ID, "password": toHash(password)}
	} else {
		return errors.New("not specify IID or ID")
	}
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(args)
	if err != nil {
		return err
	}

	c, err := result.RowsAffected()
	if c != 1 && c != 2 {
		return errors.New("failed to update password")
	}

	return err
}

func GetUser(db gfsql.DB, teamIID uint32, userID string) (*User, error) {
	stmt, err := db.PrepareNamed(SqlGetUserByID)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	user := &User{}
	err = stmt.Get(user, map[string]interface{}{"team_iid": teamIID, "user_id": userID})
	return user, err
}

func CreateUser(db gfsql.DB, teamIID uint32, userID string) (*User, error) {
	stmt, err := db.PrepareNamed(SqlCreateUser)
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
	stmt, err := db.PrepareNamed(SqlDeleteUserByIID)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"user_iid": userIID}
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

func DeleteUserByID(db gfsql.DB, teamIID uint32, userID string) error {
	stmt, err := db.PrepareNamed(SqlDeleteUserByID)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := map[string]interface{}{"team_iid": teamIID, "user_id": userID}
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

func (user *User) Delete(db gfsql.DB) error {
	if user == nil {
		return errors.New("not specify user")
	}

	if user.IID > 0 {
		return DeleteUserByIID(db, user.IID)
	} else if user.ID != "" && user.TeamIID > 0 {
		return DeleteUserByID(db, user.TeamIID, user.ID)
	}

	return errors.New("not specify user")
}
