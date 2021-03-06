package models

import (
	"testing"

	"encoding/json"

	_ "github.com/lib/pq"
)

func newJoinForTest(t *testing.T) *Join {
	return NewJoin(newDbForTest(t))
}

func TestGetHouseInvites(t *testing.T) {

	j := newJoinForTest(t)
	tearDown := TestSetup(t, j.db)
	defer tearDown(t, j.db)

	_, err := j.GetHouseInvitations(nil, 3)

	if err != nil {
		t.Errorf("Getting house invitations should work. Error: %v", err)
	}

	//fmt.Printf("result get house invite: %v", string(result))

}

func TestAddInvitation(t *testing.T) {

	j := newJoinForTest(t)
	tearDown := TestSetup(t, j.db)
	defer tearDown(t, j.db)

	var v map[string]interface{}

	inviteJSON := []byte(`{"house_id": 2, "user_id": 3}`)

	res, err := j.AddInvitation(nil, inviteJSON)

	if err != nil {
		t.Errorf("Adding invitations should work. Error: %v", err)
	}

	_ = json.Unmarshal(res, &v)
	id := v["invite_id"].(float64)

	j.DeleteRequest(nil, int64(id))
}
