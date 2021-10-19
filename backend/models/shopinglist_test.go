package models

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/util"
)

func TestGetTotalListsByOwner(t *testing.T) {
	SetupTest()

	t.Run("TestGetTotalListsByOwner", func(t *testing.T) {
		id := util.RandomInt()
		title := "title" + util.StringWithCharset(200)
		owner := "ParentListID" + util.StringWithCharset(300)
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

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
		title := "title" + util.StringWithCharset(200)
		owner := "ParentListID" + util.StringWithCharset(300)
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		err := CreateList(shoppinglist)
		if err != nil {
			t.Errorf("Error while creating Shoppinglist 1 %s", err.Error())
		}

		id2 := util.RandomInt()
		title2 := "title" + util.StringWithCharset(200)
		shoppinglist2 := Shoppinglist{
			ID:    id2,
			Title: title2,
			Owner: owner,
		}

		err = CreateList(shoppinglist2)
		if err != nil {
			t.Errorf("Error while creating Shoppinglist 2 %s", err.Error())
		}

		id3 := util.RandomInt()
		title3 := "title" + util.StringWithCharset(200)
		shoppinglist3 := Shoppinglist{
			ID:    id3,
			Title: title3,
			Owner: owner,
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
	title := "title" + util.StringWithCharset(200)
	owner := "Owner123123123123" + util.StringWithCharset(300)
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Owner: owner,
	}

	if err := CreateList(shoppinglist); err != nil {
		t.Errorf("Error while creating Shoppinglist 1 %s", err.Error())
	}

	id2 := util.RandomInt()
	title2 := "title" + util.StringWithCharset(20)
	shoppinglist2 := Shoppinglist{
		ID:    id2,
		Title: title2,
		Owner: owner,
	}

	if err := CreateList(shoppinglist2); err != nil {
		t.Errorf("Error while creating Shoppinglist 2 %s", err.Error())
	}

	lists, err := GetLists(owner, 0)
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

func TestGetListsWithOffset(t *testing.T) {
	SetupTest()

	id := util.RandomInt()
	title := "title" + util.StringWithCharset(200)
	owner := "Owner123123123123" + util.StringWithCharset(300)
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Owner: owner,
	}

	if err := CreateList(shoppinglist); err != nil {
		t.Errorf("Error while creating Shoppinglist 1 %s", err.Error())
	}

	lists, err := GetLists(owner, 1)
	if err != nil {
		t.Errorf("Error while getting the Shoppinglists %s", err.Error())
	}

	if len(lists) > 0 {
		t.Errorf("offset did not get applied")
	}

	Nil(t, err)
}

func TestBelongsShoppinglistToEmail(t *testing.T) {
	SetupTest()

	t.Run("Belongs Shoppinglist to email", func(t *testing.T) {
		id := util.RandomInt()
		title := "title" + util.StringWithCharset(200)
		owner := "Owner123123123123" + util.StringWithCharset(300)
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		if err := CreateList(shoppinglist); err != nil {
			t.Errorf("Error while creating Shoppinglist 1 %s", err.Error())
		}

		belongs, err := BelongsShoppinglistToEmail(owner, id)
		if err != nil {
			t.Errorf("Error while checking if the shoppinglist belongs to the email: %s", err)
		}

		True(t, belongs)
		Nil(t, err)
	})

	t.Run("Belongs Shoppinglist email when the shoppinglist doesn't belong to the owner", func(t *testing.T) {
		belongs, _ := BelongsShoppinglistToEmail("jdfghnkjdlfg", 0)

		False(t, belongs)
	})
}

func SetupTest() {
	Setup()
	util.Setup()
	cache.Setup(true)
}
