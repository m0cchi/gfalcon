package model_test

import (
	"github.com/m0cchi/gfalcon/model"
	"testing"
)

func TestCreateRole(t *testing.T) {
	teamID := "gfalcon"
	team, err := model.GetTeam(helper.DB, teamID)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}

	roleID := "keyaki"
	role, err := model.CreateRole(helper.DB, team.IID, roleID)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}

	if role.ID != roleID {
		t.Fatalf("expected %v but %v", roleID, role.ID)
	}

	if role.IID < 1 {
		t.Fatalf("expected IID > 0 but %v", role.IID)
	}

	if role.TeamIID != team.IID {
		t.Fatalf("expected %v but %v", team.IID, role.TeamIID)
	}
}

func TestGetRole(t *testing.T) {
	teamID := "gfalcon"
	team, err := model.GetTeam(helper.DB, teamID)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}

	roleID := "keyaki"
	role, err := model.GetRole(helper.DB, team.IID, roleID)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}

	if role.ID != roleID {
		t.Fatalf("expected %v but %v", roleID, role.ID)
	}

	if role.IID < 1 {
		t.Fatalf("expected IID > 0 but %v", role.IID)
	}

	if role.TeamIID != team.IID {
		t.Fatalf("expected %v but %v", team.IID, role.TeamIID)
	}
}

func TestDeleteRole(t *testing.T) {
	teamID := "gfalcon"
	team, err := model.GetTeam(helper.DB, teamID)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}

	roleID := "keyaki"
	role, err := model.GetRole(helper.DB, team.IID, roleID)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}

	err = role.Delete(helper.DB)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}

	_, err = model.GetRole(helper.DB, team.IID, roleID)
	if err == nil {
		t.Fatalf("missing err")
	}
}
