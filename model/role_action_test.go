package model_test

import (
	"github.com/m0cchi/gfalcon/model"
	"testing"
)

func setupRoleActionTest(t *testing.T) {
	// init test data
	var serviceIID uint32 = 1
	actionID := "keyaki"

	_, err := model.CreateAction(helper.DB, serviceIID, actionID)
	if err != nil && err != model.ErrDuplicate {
		t.Fatalf("failed create action")
	}

	teamID := "keyaki"
	team, err := model.CreateTeam(helper.DB, teamID)
	if err == model.ErrDuplicate {
		team, err = model.GetTeam(helper.DB, teamID)
		if err != nil {
			t.Fatalf("failed create team")
		}
	} else if err != nil {
		t.Fatalf("failed create team")
	}

	roleID := "keyaki"
	_, err = model.CreateRole(helper.DB, team.IID, roleID)
	if err != nil && err != model.ErrDuplicate {
		t.Fatalf("failed create team")
	}

}

func teardownRoleActionTest() {
	actionID := "keyaki"
	var serviceIID uint32 = 1
	model.DeleteActionByID(helper.DB, serviceIID, actionID)
	teamID := "keyaki"
	model.DeleteTeamByID(helper.DB, teamID)
}

func TestCreateRoleAction(t *testing.T) {
	setupRoleActionTest(t)
	var serviceIID uint32 = 1
	actionID := "keyaki"
	action, err := model.GetAction(helper.DB, serviceIID, actionID)
	if err != nil {
		t.Fatalf("test data error: missing action")
	}

	teamID := "keyaki"
	roleID := "keyaki"
	role, err := GetRole(helper.DB, teamID, roleID)

	err = model.CreateRoleAction(helper.DB, role, action)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}
}

func TestGetRoleActions(t *testing.T) {
	var serviceIID uint32 = 1
	actionID := "keyaki"
	action, err := model.GetAction(helper.DB, serviceIID, actionID)
	if err != nil {
		t.Fatalf("test data error: missing action")
	}
	teamID := "keyaki"
	roleID := "keyaki"
	role, err := GetRole(helper.DB, teamID, roleID)

	roleActions, err := model.GetRoleActions(helper.DB, role)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}

	if len(roleActions) != 1 {
		t.Fatalf("expected len == 1 but %v", len(roleActions))
	}

	if roleActions[0].RoleIID != role.IID {
		t.Fatalf("expected %v but %v", role.IID, roleActions[0].RoleIID)
	}

	if roleActions[0].ActionIID != action.IID {
		t.Fatalf("expected %v but %v", action.IID, roleActions[0].ActionIID)
	}

}

func TestDeleteRoleAction(t *testing.T) {
	var serviceIID uint32 = 1
	actionID := "keyaki"
	action, err := model.GetAction(helper.DB, serviceIID, actionID)
	if err != nil {
		t.Fatalf("test data error: missing action")
	}

	teamID := "keyaki"
	roleID := "keyaki"
	role, err := GetRole(helper.DB, teamID, roleID)

	err = model.DeleteRoleAction(helper.DB, role, action)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}
	teardownRoleActionTest()
}
