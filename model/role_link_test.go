package model_test

import (
	"github.com/m0cchi/gfalcon/model"
	"testing"
)

func setupRoleLinkTest(t *testing.T) {
	// init test data
	var serviceIID uint32 = 1
	roleID := "keyaki"

	_, err := model.CreateRole(helper.DB, serviceIID, roleID)
	if err != nil && err != model.ErrDuplicate {
		t.Fatalf("failed create role")
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

	userID := "keyaki"
	_, err = model.CreateUser(helper.DB, team.IID, userID)
	if err != nil && err != model.ErrDuplicate {
		t.Fatalf("failed create team")
	}

}

func teardownRoleLinkTest() {
	roleID := "keyaki"
	var serviceIID uint32 = 1
	model.DeleteRoleByID(helper.DB, serviceIID, roleID)
	teamID := "keyaki"
	model.DeleteTeamByID(helper.DB, teamID)
}

func TestCreateRoleList(t *testing.T) {
	setupRoleLinkTest(t)
	var serviceIID uint32 = 1
	roleID := "keyaki"
	role, err := model.GetRole(helper.DB, serviceIID, roleID)
	if err != nil {
		t.Fatalf("test data error: missing role")
	}

	teamID := "keyaki"
	userID := "keyaki"
	user, err := GetUser(helper.DB, teamID, userID)

	err = model.CreateRoleLink(helper.DB, role, user)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}
}

func TestDeleteRoleList(t *testing.T) {
	var serviceIID uint32 = 1
	roleID := "keyaki"
	role, err := model.GetRole(helper.DB, serviceIID, roleID)
	if err != nil {
		t.Fatalf("test data error: missing role")
	}

	teamID := "keyaki"
	userID := "keyaki"
	user, err := GetUser(helper.DB, teamID, userID)

	err = model.DeleteRoleLink(helper.DB, role, user)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}
	teardownRoleLinkTest()
}
