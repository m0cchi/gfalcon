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

func TestCreateTeam(t *testing.T) {
	teamID := "keyaki"
	team, err := model.CreateTeam(helper.DB, teamID)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}

	if team.ID != teamID {
		t.Fatalf("expected %v but %v", teamID, team.ID)
	}

	if team.IID < 1 {
		t.Fatalf("expected IID > 0 but %v", team.IID)
	}
}

func TestGetTeam(t *testing.T) {
	teamID := "keyaki"
	team, err := model.GetTeam(helper.DB, teamID)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}
	if team.ID != teamID {
		t.Fatalf("expected %v but %v", teamID, team.ID)
	}

	if team.IID < 1 {
		t.Fatalf("expected IID > 0 but %v", team.IID)
	}
}

func TestDeleteTeam(t *testing.T) {
	teamID := "keyaki"
	team, err := model.GetTeam(helper.DB, teamID)
	if err != nil {
		t.Fatalf("GetTeam has err: %v", err)
	}

	err = team.Delete(helper.DB)
	if err != nil {
		t.Fatalf("DeleteTeam has err: %v", err)
	}
}
