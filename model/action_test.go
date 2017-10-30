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

func TestCreateAction(t *testing.T) {
	var serviceIID uint32 = 1
	actionID := "keyaki"

	action, err := model.CreateAction(helper.DB, serviceIID, actionID)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}

	if action.ID != actionID {
		t.Fatalf("expected %v but %v", actionID, action.ID)
	}

	if action.IID < 1 {
		t.Fatalf("expected IID > 0 but %v", action.IID)
	}

	if action.ServiceIID < 1 {
		t.Fatalf("expected IID > 0 but %v", action.ServiceIID)
	}
}
