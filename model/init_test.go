package model_test

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/m0cchi/gfalcon/model"
	"os"
	"testing"
)

type Helper struct {
	DB *sqlx.DB
}

var helper Helper

func TestMain(m *testing.M) {
	datasource := os.Getenv("DATASOURCE")
	if datasource == "" {
		os.Exit(1)
	}

	db, err := sqlx.Connect("mysql", datasource)
	if err != nil {
		os.Exit(1)
	}

	defer db.Close()
	helper = Helper{db}
	code := m.Run()
	os.Exit(code)
}

func GetUser(db *sqlx.DB, teamID string, userID string) (*model.User, error) {
	team, err := model.GetTeam(db, teamID)
	if err != nil {
		return nil, err
	}
	return model.GetUser(db, team.IID, userID)
}
