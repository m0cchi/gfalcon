package model_test

import (
	"github.com/m0cchi/gfalcon/model"
	"testing"
)

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

func TestGetAction(t *testing.T) {
	var serviceIID uint32 = 1
	actionID := "keyaki"

	action, err := model.GetAction(helper.DB, serviceIID, actionID)

	if err != nil {
		t.Fatalf("has err: %v", err)
	}

	if action.ID != actionID {
		t.Fatalf("expected %v but %v", actionID, action.ID)
	}

	if action.IID < 1 {
		t.Fatalf("expected IID > 0 but %v", action.IID)
	}

	if action.ServiceIID != serviceIID {
		t.Fatalf("expected %v but %v", serviceIID, action.ServiceIID)
	}
}

func TestDeleteAction(t *testing.T) {
	var serviceIID uint32 = 1
	actionID := "keyaki"

	action, err := model.GetAction(helper.DB, serviceIID, actionID)
	err = action.Delete(helper.DB)

	if err != nil {
		t.Fatalf("has err: %v", err)
	}
	_, err = model.GetAction(helper.DB, serviceIID, actionID)

	if err == nil {
		t.Fatalf("missing err: %v", err)
	}

}
