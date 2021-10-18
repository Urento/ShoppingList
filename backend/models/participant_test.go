package models

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/pkg/util"
)

func TestGetPendingRequests(t *testing.T) {
	Setup()

	t.Run("Get Pending Requests with 1 Request Pending", func(t *testing.T) {
		id := util.RandomIntWithLength(50000)
		title := "title3332999" + util.StringWithCharset(200)
		owner := "owner999" + util.StringWithCharset(300)
		participantEmail := util.StringWithCharset(500) + "@gmail.com"
		participant := Participant{
			ParentListID: id,
			Email:        participantEmail,
			Status:       "pending",
		}
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		err := CreateList(shoppinglist)
		if err != nil {
			t.Errorf("Error while creating shoppinglist: %s", err)
		}

		_, err = AddParticipant(participant)
		if err != nil {
			t.Errorf("Error while adding participant to list: %s", err)
		}

		requests, err := GetPendingRequests(participantEmail)
		if err != nil {
			t.Errorf("Error while getting pending requests: %s", err)
		}

		Equal(t, id, requests[0].ParentListID)
		Equal(t, participantEmail, requests[0].Email)
		Equal(t, "pending", requests[0].Status)
		Nil(t, err)
	})

	t.Run("Get Pending Requests with 3 Requests", func(t *testing.T) {
		id := util.RandomIntWithLength(50000)
		title := "title3332999" + util.StringWithCharset(200)
		owner := "owner999" + util.StringWithCharset(300)
		participantEmail := util.StringWithCharset(500) + "@gmail.com"
		participantEmail2 := util.StringWithCharset(500) + "@gmail.com"
		participantEmail3 := util.StringWithCharset(500) + "@gmail.com"
		participant := Participant{
			ParentListID: id,
			Email:        participantEmail,
			Status:       "pending",
		}
		participant2 := Participant{
			ParentListID: id,
			Email:        participantEmail2,
			Status:       "pending",
		}
		participant3 := Participant{
			ParentListID: id,
			Email:        participantEmail3,
			Status:       "pending",
		}
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		err := CreateList(shoppinglist)
		if err != nil {
			t.Errorf("Error while creating shoppinglist: %s", err)
		}

		_, err = AddParticipant(participant)
		if err != nil {
			t.Errorf("Error while adding participant to list: %s", err)
		}

		_, err = AddParticipant(participant2)
		if err != nil {
			t.Errorf("Error while adding participant to list: %s", err)
		}

		_, err = AddParticipant(participant3)
		if err != nil {
			t.Errorf("Error while adding participant to list: %s", err)
		}

		requests, err := GetPendingRequests(participantEmail)
		if err != nil {
			t.Errorf("Error while getting pending requests: %s", err)
		}

		Equal(t, id, requests[0].ParentListID)
		Equal(t, participantEmail, requests[0].Email)
		Equal(t, "pending", requests[0].Status)
		Equal(t, id, requests[1].ParentListID)
		Equal(t, participantEmail, requests[1].Email)
		Equal(t, "pending", requests[1].Status)
		Equal(t, id, requests[2].ParentListID)
		Equal(t, participantEmail, requests[2].Email)
		Equal(t, "pending", requests[2].Status)
		Nil(t, err)
	})
}

func TestAcceptRequest(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(50000)
	title := "title3332999" + util.StringWithCharset(200)
	owner := "owner999" + util.StringWithCharset(300)
	participantEmail := util.StringWithCharset(500) + "@gmail.com"
	participant := Participant{
		ParentListID: id,
		Email:        participantEmail,
		Status:       "pending",
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Owner: owner,
	}

	err := CreateList(shoppinglist)
	if err != nil {
		t.Errorf("Error while creating shoppinglist: %s", err)
	}

	p, err := AddParticipant(participant)
	if err != nil {
		t.Errorf("Error while adding participant to list: %s", err)
	}

	requests, err := GetPendingRequests(participantEmail)
	if err != nil {
		t.Errorf("Error while getting pending requests: %s", err)
	}

	err = AcceptRequest(p.ID, participantEmail)
	if err != nil {
		t.Errorf("Error while accepting request: %s", err)
	}

	requestsAfter, err := GetPendingRequests(participantEmail)
	if err != nil {
		t.Errorf("Error while getting pending requests: %s", err)
	}

	if len(requestsAfter) > 0 {
		t.Errorf("request did not get accepted")
	}

	Equal(t, id, requests[0].ParentListID)
	Equal(t, participantEmail, requests[0].Email)
	Equal(t, "pending", requests[0].Status)
	Nil(t, err)
}

func TestDeleteRequest(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(50000)
	title := "title3332999" + util.StringWithCharset(200)
	owner := "owner999" + util.StringWithCharset(300)
	participantEmail := util.StringWithCharset(500) + "@gmail.com"
	participant := Participant{
		ParentListID: id,
		Email:        participantEmail,
		Status:       "pending",
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Owner: owner,
	}

	err := CreateList(shoppinglist)
	if err != nil {
		t.Errorf("Error while creating shoppinglist: %s", err)
	}

	p, err := AddParticipant(participant)
	if err != nil {
		t.Errorf("Error while adding participant to list: %s", err)
	}

	requests, err := GetPendingRequests(participantEmail)
	if err != nil {
		t.Errorf("Error while getting pending requests: %s", err)
	}

	err = DeleteRequest(p.ID, participantEmail)
	if err != nil {
		t.Errorf("Error while deleting request: %s", err)
	}

	requestsAfter, err := GetPendingRequests(participantEmail)
	if err != nil {
		t.Errorf("Error while getting pending requests: %s", err)
	}

	if len(requestsAfter) > 0 {
		t.Errorf("request did not get deleted")
	}

	Equal(t, id, requests[0].ParentListID)
	Equal(t, participantEmail, requests[0].Email)
	Equal(t, "pending", requests[0].Status)
	Nil(t, err)
}
