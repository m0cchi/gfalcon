package model_test

import (
	"github.com/m0cchi/gfalcon/model"
	"testing"
)

func setup(t *testing.T) {
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

	userID := "keyaki"
	_, err = model.CreateUser(helper.DB, team.IID, userID)
	if err != nil && err != model.ErrDuplicate {
		t.Fatalf("failed create team")
	}

}

func teardown() {
	actionID := "keyaki"
	var serviceIID uint32 = 1
	model.DeleteActionByID(helper.DB, serviceIID, actionID)
	teamID := "keyaki"
	model.DeleteTeamByID(helper.DB, teamID)
}

func TestCreateActionList(t *testing.T) {
	setup(t)
	var serviceIID uint32 = 1
	actionID := "keyaki"
	action, err := model.GetAction(helper.DB, serviceIID, actionID)
	if err != nil {
		t.Fatalf("test data error: missing action")
	}

	teamID := "keyaki"
	userID := "keyaki"
	user, err := GetUser(helper.DB, teamID, userID)

	err = model.CreateActionLink(helper.DB, action, user)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}
}

func TestGetActionList(t *testing.T) {
	var serviceIID uint32 = 1
	actionID := "keyaki"
	action, err := model.GetAction(helper.DB, serviceIID, actionID)
	if err != nil {
		t.Fatalf("test data error: missing action")
	}

	teamID := "keyaki"
	userID := "keyaki"
	user, err := GetUser(helper.DB, teamID, userID)

	actionLink, err := model.GetActionLink(helper.DB, action, user)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}

	if actionLink.ActionIID != action.IID {
		t.Fatalf("expected %v but %v", action.IID, actionLink.ActionIID)
	}

	if actionLink.UserIID != user.IID {
		t.Fatalf("expected %v but %v", user.IID, actionLink.UserIID)
	}

	if actionLink.Count != 1 {
		t.Fatalf("expected count == 1 but %v", actionLink.Count)
	}
}

func TestDeleteActionList(t *testing.T) {
	var serviceIID uint32 = 1
	actionID := "keyaki"
	action, err := model.GetAction(helper.DB, serviceIID, actionID)
	if err != nil {
		t.Fatalf("test data error: missing action")
	}

	teamID := "keyaki"
	userID := "keyaki"
	user, err := GetUser(helper.DB, teamID, userID)

	err = model.DeleteActionLink(helper.DB, action, user)
	if err != nil {
		t.Fatalf("has err: %v", err)
	}
	teardown()
}
