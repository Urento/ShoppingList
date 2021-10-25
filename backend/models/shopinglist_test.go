package models

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/util"
)

func SetupTest() {
	Setup()
	util.Setup()
	cache.Setup(false)
}

func TestGetTotalListsByOwner(t *testing.T) {
	Setup()
	util.Setup()
	cache.Setup(true)

	t.Run("TestGetTotalListsByOwner", func(t *testing.T) {
		id := util.RandomInt()
		title := "title" + util.StringWithCharset(200)
		owner := "ParentListID" + util.StringWithCharset(300)
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		err := CreateList(shoppinglist, 0, false)
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

		err := CreateList(shoppinglist, 0, false)
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

		err = CreateList(shoppinglist2, 0, false)
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

		err = CreateList(shoppinglist3, 0, false)
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
	Setup()
	util.Setup()
	cache.Setup(true)

	id := util.RandomInt()
	title := "title" + util.StringWithCharset(200)
	owner := "Owner123123123123" + util.StringWithCharset(300)
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Owner: owner,
	}

	if err := CreateList(shoppinglist, 0, false); err != nil {
		t.Errorf("Error while creating Shoppinglist 1 %s", err.Error())
	}

	id2 := util.RandomInt()
	title2 := "title" + util.StringWithCharset(20)
	shoppinglist2 := Shoppinglist{
		ID:    id2,
		Title: title2,
		Owner: owner,
	}

	if err := CreateList(shoppinglist2, 0, false); err != nil {
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
	Equal(t, title, lists[0].Title)
	Equal(t, id, lists[0].ID)
	Equal(t, owner, lists[1].Owner)
	Equal(t, title2, lists[1].Title)
	Equal(t, id2, lists[1].ID)
}

func TestGetListsWithOffset(t *testing.T) {
	Setup()
	util.Setup()
	cache.Setup(true)

	id := util.RandomInt()
	title := "title" + util.StringWithCharset(200)
	owner := "Owner123123123123" + util.StringWithCharset(300)
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Owner: owner,
	}

	if err := CreateList(shoppinglist, 0, false); err != nil {
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
	Setup()
	util.Setup()
	cache.Setup(true)

	t.Run("Belongs Shoppinglist to email", func(t *testing.T) {
		id := util.RandomInt()
		title := "title" + util.StringWithCharset(200)
		owner := "Owner123123123123" + util.StringWithCharset(300)
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		if err := CreateList(shoppinglist, 0, false); err != nil {
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

func TestCreate(t *testing.T) {
	Setup()
	util.Setup()
	cache.Setup(true)

	t.Run("Create and Check", func(t *testing.T) {
		id := util.RandomIntWithLength(5000)
		title := "title" + util.StringWithCharset(200)
		owner := "owner" + util.StringWithCharset(300)
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		err := CreateList(shoppinglist, 0, false)
		if err != nil {
			t.Errorf("Failed to create shoppinglist %s", err.Error())
		}

		exists, err := ExistByID(id)
		if err != nil || !exists {
			t.Errorf("Shoppinglist did not get created %s", err.Error())
		}

		l, err := GetList(id, owner)
		if err != nil {
			t.Errorf("Error while getting the shoppinglist: %s", err.Error())
		}

		err = DeleteList(id)
		if err != nil {
			t.Errorf("Shoppinglist couldn't be deleted %s", err.Error())
		}

		Equal(t, id, l.ID)
		Equal(t, title, l.Title)
		Equal(t, owner, l.Owner)
		Equal(t, nil, err)
	})

	t.Run("Create and Edit", func(t *testing.T) {
		id := util.RandomIntWithLength(5000)
		title := "title33232999" + util.StringWithCharset(200)
		owner := "999owner999" + util.StringWithCharset(300)
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		err := CreateList(shoppinglist, 0, false)
		if err != nil {
			t.Errorf("Failed to create shoppinglist %s", err.Error())
		}

		title2 := util.StringWithCharset(9) + util.StringWithCharset(200)
		shoppinglist = Shoppinglist{
			ID:    id,
			Title: title2,
			Owner: owner,
		}

		err = EditList(id, shoppinglist)
		if err != nil {
			t.Errorf("Failed to edit shoppinglist %s", err.Error())
		}

		exists, err := ExistByID(id)
		if err != nil || !exists {
			t.Errorf("Shoppinglist did not get created %s", err.Error())
		}

		l, err := GetList(id, owner)
		if err != nil {
			t.Errorf("Shoppinglist not found %s", err.Error())
		}

		err = DeleteList(id)
		if err != nil {
			t.Errorf("Shoppinglist couldn't be deleted %s", err.Error())
		}

		Equal(t, id, l.ID)
		Equal(t, title2, l.Title)
		Equal(t, owner, l.Owner)
		Equal(t, nil, err)
		NotEqual(t, title, l.Title)
	})
}

func TestExistsByID(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(500000)
	title := "titlesdfgdsghdshgfdzhjf" + util.StringWithCharset(20000)
	owner := "ownersthfdghdfhfdthfxgdh" + util.StringWithCharset(30000)
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Owner: owner,
	}

	err := CreateList(shoppinglist, 0, false)
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	exists, err := ExistByID(id)
	if err != nil || !exists {
		t.Errorf("Shoppinglist did not get created %s", err.Error())
	}

	err = DeleteList(id)
	if err != nil {
		t.Errorf("Shoppinglist couldn't be deleted %s", err.Error())
	}

	Equal(t, true, exists)
	Equal(t, nil, err)
}

func TestAddItem(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(50000)
	itemID := util.RandomIntWithLength(50000)
	title := "title3332999" + util.StringWithCharset(20000)
	owner := "owner999" + util.StringWithCharset(30000)
	items := []*Item{
		{
			ParentListID: id,
			ItemID:       itemID,
			Title:        util.StringWithCharset(10000),
			Position:     util.RandomPosition(),
		},
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Items: items,
		Owner: owner,
	}

	err := CreateList(shoppinglist, 0, false)
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	item, err := AddItem(*items[0])
	if err != nil {
		t.Errorf("Failed to edit shoppinglist %s", err.Error())
	}

	Equal(t, nil, err)
	NotNil(t, item)
}

func TestGetList(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(50000)
	title := "titlesdfgdsghdshgfdzhjf" + util.StringWithCharset(20000)
	owner := "ownersthfdghdfhfdthfxgdh" + util.StringWithCharset(3000)
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Owner: owner,
	}

	err := CreateList(shoppinglist, 0, false)
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	list, err := GetList(id, owner)
	if err != nil {
		t.Errorf("Error while getting list: %s", err)
	}

	Equal(t, id, list.ID)
	Equal(t, title, list.Title)
}

func TestGetItems(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(50000)
	itemID := util.RandomIntWithLength(50000)
	title := "title3332999" + util.StringWithCharset(20000)
	owner := "owner999" + util.StringWithCharset(30000)
	items := []*Item{
		{
			ParentListID: id,
			ItemID:       itemID,
			Title:        util.StringWithCharset(10000),
			Position:     util.RandomPosition(),
		},
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Items: items,
		Owner: owner,
	}

	err := CreateList(shoppinglist, 0, false)
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	_, err = AddItem(*items[0])
	if err != nil {
		t.Errorf("Error while adding item: %s", err)
	}

	itemsInList, _ := GetItems(id)

	Equal(t, items[0].Title, itemsInList[0].Title)
	Equal(t, items[0].Bought, itemsInList[0].Bought)
	Equal(t, items[0].ParentListID, itemsInList[0].ParentListID)
	Equal(t, items[0].Position, itemsInList[0].Position)
	Equal(t, items[0].ItemID, itemsInList[0].ItemID)
	Equal(t, nil, err)
}

func TestGetLastPosition(t *testing.T) {
	Setup()

	t.Run("Get last position with two items", func(t *testing.T) {
		id := util.RandomIntWithLength(7000)
		itemID := util.RandomIntWithLength(50000)
		title := "title3332999" + util.StringWithCharset(20000)
		owner := "owner999" + util.StringWithCharset(300)
		position := util.RandomPosition()
		items := []*Item{
			{
				ParentListID: id,
				ItemID:       itemID,
				Title:        util.StringWithCharset(100),
				Position:     position,
			},
		}
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Items: items,
			Owner: owner,
		}

		err := CreateList(shoppinglist, 0, false)
		if err != nil {
			t.Errorf("Failed to create shoppinglist %s", err.Error())
		}

		_, err = AddItem(*items[0])
		if err != nil {
			t.Errorf("Error while adding item: %s", err)
		}

		position2 := util.RandomPosition()
		items2 := []*Item{
			{
				ParentListID: id,
				ItemID:       itemID,
				Title:        util.StringWithCharset(100),
				Position:     position2,
			},
		}

		_, err = AddItem(*items2[0])
		if err != nil {
			t.Errorf("Error while adding item: %s", err)
		}

		lastPosition, err := GetLastPosition(id)
		if err != nil {
			t.Errorf("Error while getting last position: %s", err)
		}

		if lastPosition < position || lastPosition < position2 {
			t.Errorf("Last Position is not the largest number")
		}
	})

	t.Run("Get Last Position without any items", func(t *testing.T) {
		id := util.RandomIntWithLength(5000)
		title := "title3332999" + util.StringWithCharset(200)
		owner := "owner999" + util.StringWithCharset(300)
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		err := CreateList(shoppinglist, 0, false)
		if err != nil {
			t.Errorf("Failed to create shoppinglist %s", err.Error())
		}

		lastPosition, _ := GetLastPosition(id)

		Equal(t, int64(0), lastPosition)
	})
}

func TestGetItem(t *testing.T) {
	Setup()

	t.Run("Get Items", func(t *testing.T) {
		id := util.RandomIntWithLength(50000)
		itemID := util.RandomIntWithLength(200)
		title := "title3332999" + util.StringWithCharset(200)
		owner := "owner999" + util.StringWithCharset(300)
		position := util.RandomPosition()
		itemTitle := util.StringWithCharset(100)
		items := []*Item{
			{
				ParentListID: id,
				ItemID:       itemID,
				Title:        itemTitle,
				Position:     position,
			},
		}
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Items: items,
			Owner: owner,
		}

		err := CreateList(shoppinglist, 0, false)
		if err != nil {
			t.Errorf("Failed to create shoppinglist %s", err.Error())
		}

		_, err = AddItem(*items[0])
		if err != nil {
			t.Errorf("Error while adding item: %s", err)
		}

		i, err := GetItem(id, itemID)
		if err != nil {
			t.Errorf("Error while getting item: %s", err)
		}

		Equal(t, id, i.ParentListID)
		Equal(t, itemID, i.ItemID)
		Equal(t, itemTitle, i.Title)
		Equal(t, position, i.Position)
	})

	t.Run("Get Item that doesn't exist", func(t *testing.T) {
		id := util.RandomIntWithLength(500)
		itemID := util.RandomIntWithLength(200)

		i, err := GetItem(id, itemID)
		t.Log(i)

		NotEqual(t, nil, err)
	})
}

func TestUpdateItem(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(50000)
	itemID := util.RandomIntWithLength(50000)
	title := "title3332999" + util.StringWithCharset(00020)
	owner := "owner999" + util.StringWithCharset(30000)
	position := util.RandomPosition()
	itemTitle := util.StringWithCharset(10000)
	items := []*Item{
		{
			ParentListID: id,
			ItemID:       itemID,
			Title:        itemTitle,
			Position:     position,
		},
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Items: items,
		Owner: owner,
	}

	err := CreateList(shoppinglist, 0, false)
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	_, err = AddItem(*items[0])
	if err != nil {
		t.Errorf("Error while adding item: %s", err)
	}

	i, err := GetItem(id, itemID)
	if err != nil {
		t.Errorf("Error while getting item: %s", err)
	}

	newTitle := util.StringWithCharset(10000)
	newPosition := util.RandomPosition()
	newBought := util.RandomBoolean()
	updatedItem := Item{
		ItemID:       itemID,
		ParentListID: id,
		Title:        newTitle,
		Position:     newPosition,
		Bought:       newBought,
	}

	err = UpdateItem(updatedItem)
	if err != nil {
		t.Errorf("Error while updating item: %s", err)
	}

	i2, err := GetItem(id, itemID)
	if err != nil {
		t.Errorf("Error while getting item: %s", err)
	}

	Equal(t, id, i.ParentListID)
	Equal(t, itemID, i.ItemID)
	Equal(t, itemTitle, i.Title)
	Equal(t, position, i.Position)
	Equal(t, newTitle, i2.Title)
	Equal(t, newPosition, i2.Position)
	Equal(t, id, i2.ParentListID)
	Equal(t, itemID, i2.ItemID)
	Equal(t, newBought, i2.Bought)
}

func TestAddParticipant(t *testing.T) {
	Setup()

	t.Run("Add Participant", func(t *testing.T) {
		id := util.RandomInt()
		title := "title3332999" + util.StringWithCharset(20000)
		owner := "owner999" + util.StringWithCharset(30000)
		participantEmail := util.RandomEmail()
		participantStatus := util.StringWithCharset(20000)
		participant := Participant{
			ParentListID: id,
			Email:        participantEmail,
			Status:       participantStatus,
		}
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		err := CreateList(shoppinglist, 0, false)
		if err != nil {
			t.Errorf("Failed to create shoppinglist %s", err.Error())
		}

		_, err = AddParticipant(participant)
		if err != nil {
			t.Errorf("Error while adding participant: %s", err)
		}

		Nil(t, err)
	})

	t.Run("Add Participant when Shoppinglist doesn't exist", func(t *testing.T) {
		id := util.RandomIntWithLength(50000)
		participantEmail := util.RandomEmail()
		participantStatus := util.StringWithCharset(20000)
		participant := Participant{
			ParentListID: id,
			Email:        participantEmail,
			Status:       participantStatus,
		}

		_, err := AddParticipant(participant)

		NotNil(t, err)
	})
}

func TestGetParticipants(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(50000)
	title := "title3332999" + util.StringWithCharset(20000)
	owner := "owner999" + util.StringWithCharset(30000)
	participantEmail := util.RandomEmail()
	participantStatus := util.StringWithCharset(20000)
	participant := Participant{
		ParentListID: id,
		Email:        participantEmail,
		Status:       participantStatus,
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Owner: owner,
	}

	err := CreateList(shoppinglist, 0, false)
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	_, err = AddParticipant(participant)
	if err != nil {
		t.Errorf("Error while adding participant: %s", err)
	}

	participants, err := GetParticipants(id)
	if err != nil {
		t.Errorf("Error while getting participants: %s", err)
	}

	Equal(t, participantEmail, participants[0].Email)
	Equal(t, participantStatus, participants[0].Status)
	Equal(t, id, participants[0].ParentListID)
}

func TestRemoveParticipant(t *testing.T) {
	Setup()

	t.Run("Remove Participant", func(t *testing.T) {
		id := util.RandomInt() + 50000
		title := "title3332999" + util.StringWithCharset(20000)
		owner := "owner999" + util.StringWithCharset(30000)
		participantEmail := util.RandomEmail()
		participantStatus := util.StringWithCharset(20000)
		participant := Participant{
			ParentListID: id,
			Email:        participantEmail,
			Status:       participantStatus,
		}
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		err := CreateList(shoppinglist, 0, false)
		if err != nil {
			t.Errorf("Failed to create shoppinglist %s", err)
		}

		_, err = AddParticipant(participant)
		if err != nil {
			t.Errorf("Error while adding participant: %s", err)
		}

		participants, err := GetParticipants(id)
		if err != nil {
			t.Errorf("Error while getting participants: %s", err)
		}

		delErr := RemoveParticipant(id, participants[0].ID)
		if delErr != nil {
			t.Errorf("Error while removing participant: %s", err)
		}

		participantsAfter, err := GetParticipants(id)
		if err != nil {
			t.Errorf("Error while getting participants: %s", err)
		}

		if len(participantsAfter) > 0 {
			t.Errorf("Participant did not get deleted")
		}

		Nil(t, err)
	})

	t.Run("Remove Participant when Shoppinglist doesn't exist", func(t *testing.T) {
		id := util.RandomIntWithLength(50000)

		err := RemoveParticipant(id, 348756324753645)

		NotNil(t, err)
	})
}

func TestDeleteItem(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(50000)
	itemID := util.RandomIntWithLength(50000)
	title := "title3332999" + util.StringWithCharset(20000)
	owner := "owner999" + util.StringWithCharset(30000)
	items := []*Item{
		{
			ParentListID: id,
			ItemID:       itemID,
			Title:        util.StringWithCharset(10000),
			Position:     util.RandomPosition(),
		},
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Items: items,
		Owner: owner,
	}

	err := CreateList(shoppinglist, 0, false)
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	item, err := AddItem(*items[0])
	if err != nil {
		t.Errorf("Failed to edit shoppinglist %s", err.Error())
	}

	err = DeleteItem(id, itemID)
	if err != nil {
		t.Errorf("Error while deleting item: %s", err)
	}

	it, err := GetItems(id)
	if err != nil {
		t.Errorf("Error while getting items: %s", err)
	}
	t.Log(it)

	if len(it) > 0 {
		t.Errorf("Item didn't get deleted")
	}

	Equal(t, nil, err)
	NotNil(t, item)
}

func TestUpdateItems(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(50000)
	itemID := util.RandomIntWithLength(50000)
	itemID2 := util.RandomIntWithLength(50000)
	title := "title3332999" + util.StringWithCharset(200)
	owner := "owner999" + util.StringWithCharset(300)
	items := []*Item{
		{
			ParentListID: id,
			ItemID:       itemID,
			Title:        util.StringWithCharset(100),
			Position:     util.RandomPosition(),
		},
		{
			ParentListID: id,
			ItemID:       itemID2,
			Title:        util.StringWithCharset(100),
			Position:     util.RandomPosition(),
		},
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

	for _, item := range items {
		_, err = AddItem(*item)
		if err != nil {
			t.Errorf("Error while adding item: %s", err)
		}
	}

	title3 := util.StringWithCharset(100)
	position3 := util.RandomPosition()
	title4 := util.StringWithCharset(100)
	position4 := util.RandomPosition()
	updateItems := []Item{
		{
			ParentListID: id,
			ItemID:       itemID,
			Title:        title3,
			Position:     position3,
		},
		{
			ParentListID: id,
			ItemID:       itemID2,
			Title:        title4,
			Position:     position4,
		},
	}

	err = UpdateItems(id, updateItems)
	if err != nil {
		t.Errorf("Error while updating items: %s", err)
	}

	i, err := GetItems(id)
	if err != nil {
		t.Errorf("Error while getting items: %s", err)
	}

	Equal(t, title3, i[0].Title)
	Equal(t, position3, i[0].Position)
	Equal(t, title4, i[1].Title)
	Equal(t, position4, i[1].Position)
}

func TestGetListsByParticipant(t *testing.T) {
	Setup()

	t.Run("Get Lists By Participant with 1 list", func(t *testing.T) {
		id := util.RandomIntWithLength(50000)
		title := "title3332999" + util.StringWithCharset(200)
		owner := "owner999" + util.StringWithCharset(300)
		participantEmail := util.RandomEmail()
		participant := Participant{
			ParentListID: id,
			Email:        participantEmail,
			Status:       "accepted",
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

		lists, err := GetListsByParticipant(participantEmail)
		if err != nil {
			t.Errorf("Error while getting lists by participant: %s", err)
		}

		Equal(t, owner, lists[0].Owner)
		Equal(t, title, lists[0].Title)
		Equal(t, id, lists[0].ID)
	})

	t.Run("Get Lists By Participant with 3 lists", func(t *testing.T) {
		id := util.RandomIntWithLength(90000)
		title := "title3332999" + util.StringWithCharset(200)
		owner := "owner999" + util.StringWithCharset(300)
		participantEmail := util.RandomEmail()
		participantEmail2 := util.RandomEmail()
		participantEmail3 := util.RandomEmail()
		participant := Participant{
			ParentListID: id,
			Email:        participantEmail,
			Status:       "accepted",
			RequestFrom:  owner,
		}
		participant2 := Participant{
			ParentListID: id,
			Email:        participantEmail2,
			Status:       "accepted",
			RequestFrom:  owner,
		}
		participant3 := Participant{
			ParentListID: id,
			Email:        participantEmail3,
			Status:       "accepted",
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

		lists, err := GetListsByParticipant(participantEmail)
		if err != nil {
			t.Errorf("Error while getting lists by participant: %s", err)
		}

		Equal(t, owner, lists[0].Owner)
		Equal(t, title, lists[0].Title)
		Equal(t, id, lists[0].ID)
		Equal(t, owner, lists[1].Owner)
		Equal(t, title, lists[1].Title)
		Equal(t, id, lists[1].ID)
		Equal(t, owner, lists[2].Owner)
		Equal(t, title, lists[2].Title)
		Equal(t, id, lists[2].ID)
	})
}

func TestGetListWithoutOwner(t *testing.T) {
	SetupTestAuth()

	id := util.RandomIntWithLength(90000)
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

	list, err := GetListWithoutOwner(id)
	if err != nil {
		t.Errorf("Error while getting list without owner: %s", err)
	}

	Equal(t, id, list.ID)
	Equal(t, title, list.Title)
	Equal(t, owner, list.Owner)
	Nil(t, err)
}
