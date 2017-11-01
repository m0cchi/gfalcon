package model_test

import (
	"github.com/m0cchi/gfalcon/model"
	"testing"
)

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
