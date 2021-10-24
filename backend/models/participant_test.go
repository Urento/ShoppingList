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
			RequestFrom:  owner,
		}
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		err := CreateList(shoppinglist, 0, false)
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
		Equal(t, owner, requests[0].RequestFrom)
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
			RequestFrom:  owner,
		}
		participant2 := Participant{
			ParentListID: id,
			Email:        participantEmail2,
			Status:       "pending",
			RequestFrom:  owner,
		}
		participant3 := Participant{
			ParentListID: id,
			Email:        participantEmail3,
			Status:       "pending",
			RequestFrom:  owner,
		}
		participant4 := Participant{
			ParentListID: id,
			Email:        participantEmail3,
			Status:       "pending",
			RequestFrom:  owner,
		}
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		err := CreateList(shoppinglist, 0, false)
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

		_, err = AddParticipant(participant4)
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
		Equal(t, owner, requests[0].RequestFrom)
		Equal(t, id, requests[1].ParentListID)
		Equal(t, participantEmail, requests[1].Email)
		Equal(t, "pending", requests[1].Status)
		Equal(t, owner, requests[1].RequestFrom)
		Equal(t, id, requests[2].ParentListID)
		Equal(t, participantEmail, requests[2].Email)
		Equal(t, "pending", requests[2].Status)
		Equal(t, owner, requests[2].RequestFrom)
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
		RequestFrom:  owner,
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Owner: owner,
	}

	err := CreateList(shoppinglist, 0, false)
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
	Equal(t, owner, requests[0].RequestFrom)
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
		RequestFrom:  owner,
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Owner: owner,
	}

	err := CreateList(shoppinglist, 0, false)
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
	Equal(t, owner, requests[0].RequestFrom)
	Nil(t, err)
}

func TestGetPendingRequestsFromShoppinglist(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(50000)
	title := "title3332999" + util.StringWithCharset(200)
	owner := "owner999" + util.StringWithCharset(300)
	participantEmail := util.StringWithCharset(500) + "@gmail.com"
	participant := Participant{
		ParentListID: id,
		Email:        participantEmail,
		Status:       "pending",
		RequestFrom:  owner,
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Owner: owner,
	}

	err := CreateList(shoppinglist, 0, false)
	if err != nil {
		t.Errorf("Error while creating shoppinglist: %s", err)
	}

	_, err = AddParticipant(participant)
	if err != nil {
		t.Errorf("Error while adding participant to list: %s", err)
	}

	requests, err := GetPendingRequestsFromShoppinglist(owner, id)
	if err != nil {
		t.Errorf("Error while getting pending requess from shoppinglist: %s", err)
	}

	Equal(t, id, requests[0].ParentListID)
	Equal(t, "pending", requests[0].Status)
	Equal(t, participantEmail, requests[0].Email)
	Equal(t, owner, requests[0].RequestFrom)
}

func TestIsParticipantAlreadyIncluded(t *testing.T) {
	Setup()

	t.Run("Is Participant already included", func(t *testing.T) {
		id := util.RandomInt()
		title := "title3332999" + util.StringWithCharset(200)
		owner := "owner999" + util.StringWithCharset(300)
		participantEmail := util.RandomEmail()
		participant := Participant{
			ParentListID: id,
			Email:        participantEmail,
			Status:       "pending",
			RequestFrom:  owner,
		}
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		err := CreateList(shoppinglist, 0, false)
		if err != nil {
			t.Errorf("Error while creating shoppinglist: %s", err)
		}

		_, err = AddParticipant(participant)
		if err != nil {
			t.Errorf("Error while adding participant to list: %s", err)
		}

		included, err := IsParticipantAlreadyIncluded(participantEmail, id)
		if err != nil {
			t.Errorf("Error while checking if the participant is already included: %s", err)
		}

		True(t, included)
	})

	t.Run("Is Participant already included but the participant isn't included", func(t *testing.T) {
		id := util.RandomIntWithLength(9000000)
		title := "title3332999" + util.StringWithCharset(200)
		owner := "owner999" + util.StringWithCharset(300)
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		err := CreateList(shoppinglist, 0, false)
		if err != nil {
			t.Errorf("Error while creating shoppinglist: %s", err)
		}

		included, err := IsParticipantAlreadyIncluded("jndfgh", id)
		if err != nil {
			t.Errorf("Error while checking if the participant is already included: %s", err)
		}

		False(t, included)
	})
}

func TestDeleteAll(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(50000)
	title := "title3332999" + util.StringWithCharset(200)
	owner := "owner999" + util.StringWithCharset(300)
	participantEmail := util.StringWithCharset(500) + "@gmail.com"
	participantEmail2 := util.StringWithCharset(500) + "@gmail.com"
	participantEmail3 := util.StringWithCharset(500) + "@gmail.com"
	participantEmail4 := util.StringWithCharset(500) + "@gmail.com"
	participant := Participant{
		ParentListID: id,
		Email:        participantEmail,
		Status:       "pending",
		RequestFrom:  owner,
	}
	participant2 := Participant{
		ParentListID: id,
		Email:        participantEmail2,
		Status:       "pending",
		RequestFrom:  owner,
	}
	participant3 := Participant{
		ParentListID: id,
		Email:        participantEmail3,
		Status:       "pending",
		RequestFrom:  owner,
	}
	participant4 := Participant{
		ParentListID: id,
		Email:        participantEmail4,
		Status:       "pending",
		RequestFrom:  owner,
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Owner: owner,
	}

	err := CreateList(shoppinglist, 0, false)
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

	_, err = AddParticipant(participant4)
	if err != nil {
		t.Errorf("Error while adding participant to list: %s", err)
	}

	err = DeleteAll(owner)
	if err != nil {
		t.Errorf("Error while deleting all requests: %s", err)
	}

	requests, err := GetPendingRequests(owner)
	if err != nil {
		t.Errorf("Error while getting pending requests: %s", err)
	}

	if len(requests) > 0 {
		t.Errorf("Not all Requests got deleted")
	}

	Nil(t, err)
}

func TestLeaveShoppinglist(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(50000)
	title := "title3332999" + util.StringWithCharset(200)
	owner := "owner999" + util.StringWithCharset(300)
	participantEmail := util.StringWithCharset(500) + "@gmail.com"
	participant := Participant{
		ParentListID: id,
		Email:        participantEmail,
		Status:       "pending",
		RequestFrom:  owner,
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Owner: owner,
	}

	err := CreateList(shoppinglist, 0, false)
	if err != nil {
		t.Errorf("Error while creating shoppinglist: %s", err)
	}

	p, err := AddParticipant(participant)
	if err != nil {
		t.Errorf("Error while adding participant to list: %s", err)
	}

	err = AcceptRequest(p.ID, p.Email)
	if err != nil {
		t.Errorf("Error while accepting request: %s", err)
	}

	err = LeaveShoppinglist(p.ID, p.Email)
	if err != nil {
		t.Errorf("Error while leaving shoppinglist: %s", err)
	}

	requests, err := GetPendingRequests(owner)
	if err != nil {
		t.Errorf("Error while getting pending requests: %s", err)
	}

	if len(requests) > 0 {
		t.Errorf("Not all Requests got deleted")
	}

	Nil(t, err)
}
