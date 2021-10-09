package models

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/pkg/util"
)

func TestGetTotalListsByOwner(t *testing.T) {
	SetupTest()

	t.Run("TestGetTotalListsByOwner", func(t *testing.T) {
		id := util.RandomInt()
		title := "title" + util.StringWithCharset(20)
		owner := "ParentListID" + util.StringWithCharset(30)
		participants := []*Participant{
			{
				ParentListID: id,
				Email:        util.RandomEmail(),
			},
			{
				ParentListID: id,
				Email:        util.RandomEmail(),
			},
		}
		shoppinglist := map[string]interface{}{
			"id":           id,
			"title":        title,
			"owner":        owner,
			"participants": participants,
		}
		t.Log(shoppinglist)

		err := CreateList(shoppinglist)
		if err != nil {
			t.Errorf("Error while creating Shoppinglist %s", err.Error())
		}

		count, err := GetTotalListsByOwner(owner)
		if err != nil {
			t.Errorf("Error while getting the total lists by ParentListID %s", err.Error())
		}

		Equal(t, count, int64(1))
	})

	t.Run("TestGetTotalListsByOwnerWithMultipleLists", func(t *testing.T) {
		id := util.RandomInt()
		title := "title" + util.StringWithCharset(20)
		owner := "ParentListID" + util.StringWithCharset(30)
		participants := []*Participant{
			{
				ParentListID: id,
				Email:        util.RandomEmail(),
			},
			{
				ParentListID: id,
				Email:        util.RandomEmail(),
			},
		}
		shoppinglist := map[string]interface{}{
			"id":           id,
			"title":        title,
			"owner":        owner,
			"participants": participants,
		}

		err := CreateList(shoppinglist)
		if err != nil {
			t.Errorf("Error while creating Shoppinglist 1 %s", err.Error())
		}

		id2 := util.RandomInt()
		title2 := "title" + util.StringWithCharset(20)
		participants2 := []*Participant{
			{
				ParentListID: id,
				Email:        util.RandomEmail(),
			},
			{
				ParentListID: id,
				Email:        util.RandomEmail(),
			},
		}
		shoppinglist2 := map[string]interface{}{
			"id":           id2,
			"title":        title2,
			"owner":        owner,
			"participants": participants2,
		}

		err = CreateList(shoppinglist2)
		if err != nil {
			t.Errorf("Error while creating Shoppinglist 2 %s", err.Error())
		}

		id3 := util.RandomInt()
		title3 := "title" + util.StringWithCharset(20)
		participants3 := []*Participant{
			{
				ParentListID: id,
				Email:        util.RandomEmail(),
			},
			{
				ParentListID: id,
				Email:        util.RandomEmail(),
			},
		}
		shoppinglist3 := map[string]interface{}{
			"id":           id3,
			"title":        title3,
			"owner":        owner,
			"participants": participants3,
		}

		err = CreateList(shoppinglist3)
		if err != nil {
			t.Errorf("Error while creating Shoppinglist 3 %s", err.Error())
		}

		count, err := GetTotalListsByOwner(owner)
		if err != nil {
			t.Errorf("Error while getting the total lists by ParentListID %s", err.Error())
		}

		Equal(t, count, int64(3))
	})
}

func TestGetListsByOwner(t *testing.T) {
	SetupTest()

	id := util.RandomInt()
	title := "title" + util.StringWithCharset(20)
	owner := "Owner123123123123" + util.StringWithCharset(30)
	participants := []*Participant{
		{
			ParentListID: id,
			Email:        util.RandomEmail(),
		},
		{
			ParentListID: id,
			Email:        util.RandomEmail(),
		},
	}
	shoppinglist := map[string]interface{}{
		"id":           id,
		"title":        title,
		"owner":        owner,
		"participants": participants,
	}

	if err := CreateList(shoppinglist); err != nil {
		t.Errorf("Error while creating Shoppinglist 1 %s", err.Error())
	}

	id2 := util.RandomInt()
	title2 := "title" + util.StringWithCharset(20)
	participants2 := []*Participant{
		{
			ParentListID: id,
			Email:        util.RandomEmail(),
		},
		{
			ParentListID: id,
			Email:        util.RandomEmail(),
		},
	}
	shoppinglist2 := map[string]interface{}{
		"id":           id2,
		"title":        title2,
		"owner":        owner,
		"participants": participants2,
	}

	if err := CreateList(shoppinglist2); err != nil {
		t.Errorf("Error while creating Shoppinglist 2 %s", err.Error())
	}

	lists, err := GetLists(owner)
	if err != nil {
		t.Errorf("Error while getting the Shoppinglists %s", err.Error())
	}

	if len(lists) <= 0 {
		t.Errorf("List is Empty")
	}

	Equal(t, owner, lists[0].Owner)
	Equal(t, owner, lists[1].Owner)
	Equal(t, title, lists[0].Title)
	Equal(t, id, lists[0].ID)
	Equal(t, title2, lists[1].Title)
	Equal(t, id2, lists[1].ID)
}

func SetupTest() {
	Setup()
	util.Setup()
}
