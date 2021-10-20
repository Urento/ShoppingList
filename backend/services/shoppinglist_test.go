package services

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/models"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/util"
)

func Setup() {
	models.Setup()
	util.Setup()
	cache.Setup(true)
}

func TestCreate(t *testing.T) {
	Setup()

	t.Run("Create and Check", func(t *testing.T) {
		id := util.RandomIntWithLength(5000)
		title := "title" + util.StringWithCharset(200)
		owner := "owner" + util.StringWithCharset(300)
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		created, err := shoppinglist.Create(0, false)
		if err != nil || !created {
			t.Errorf("Failed to create shoppinglist %s", err.Error())
		}

		exists, err := shoppinglist.ExistsByID()
		if err != nil || !exists {
			t.Errorf("Shoppinglist did not get created %s", err.Error())
		}

		l, err := shoppinglist.GetList()
		if err != nil {
			t.Errorf("Error while getting the shoppinglist: %s", err.Error())
		}

		err = shoppinglist.Delete(0, false)
		if err != nil {
			t.Errorf("Shoppinglist couldn't be deleted %s", err.Error())
		}

		Equal(t, id, l.ID)
		Equal(t, title, l.Title)
		Equal(t, owner, l.Owner)
		Equal(t, nil, err)
		Equal(t, true, created)
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

		created, err := shoppinglist.Create(0, false)
		if err != nil {
			t.Errorf("Failed to create shoppinglist %s", err.Error())
		}

		title2 := util.StringWithCharset(9) + util.StringWithCharset(200)
		shoppinglist = Shoppinglist{
			ID:    id,
			Title: title2,
			Owner: owner,
		}

		err = shoppinglist.Edit(0, false)
		if err != nil {
			t.Errorf("Failed to edit shoppinglist %s", err.Error())
		}

		exists, err := shoppinglist.ExistsByID()
		if err != nil || !exists {
			t.Errorf("Shoppinglist did not get created %s", err.Error())
		}

		l, err := shoppinglist.GetList()
		if err != nil {
			t.Errorf("Shoppinglist not found %s", err.Error())
		}

		err = shoppinglist.Delete(0, false)
		if err != nil {
			t.Errorf("Shoppinglist couldn't be deleted %s", err.Error())
		}

		Equal(t, true, created)
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

	created, err := shoppinglist.Create(0, false)
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	list := Shoppinglist{ID: id}
	exists, err := list.ExistsByID()
	if err != nil || !exists {
		t.Errorf("Shoppinglist did not get created %s", err.Error())
	}

	err = list.Delete(0, false)
	if err != nil {
		t.Errorf("Shoppinglist couldn't be deleted %s", err.Error())
	}

	Equal(t, true, created)
	Equal(t, true, exists)
	Equal(t, nil, err)
}

func TestAddItem(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(50000)
	itemID := util.RandomIntWithLength(50000)
	title := "title3332999" + util.StringWithCharset(20000)
	owner := "owner999" + util.StringWithCharset(30000)
	items := &models.Item{
		ParentListID: id,
		ItemID:       itemID,
		Title:        util.StringWithCharset(10000),
		Position:     util.RandomPosition(),
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Items: *items,
		Owner: owner,
	}

	created, err := shoppinglist.Create(0, false)
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	item, err := shoppinglist.AddItem()
	if err != nil {
		t.Errorf("Failed to edit shoppinglist %s", err.Error())
	}

	Equal(t, true, created)
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

	created, err := shoppinglist.Create(0, false)
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	list, err := shoppinglist.GetList()
	if err != nil {
		t.Errorf("Error while getting list: %s", err)
	}

	Equal(t, true, created)
	Equal(t, id, list.ID)
	Equal(t, title, list.Title)
}

func TestGetItems(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(50000)
	itemID := util.RandomIntWithLength(50000)
	title := "title3332999" + util.StringWithCharset(20000)
	owner := "owner999" + util.StringWithCharset(30000)
	items := &models.Item{
		ParentListID: id,
		ItemID:       itemID,
		Title:        util.StringWithCharset(10000),
		Position:     util.RandomPosition(),
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Items: *items,
		Owner: owner,
	}

	created, err := shoppinglist.Create(0, false)
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	_, err = shoppinglist.AddItem()
	if err != nil {
		t.Errorf("Error while adding item: %s", err)
	}

	if !created {
		t.Errorf("Error while creating shoppinglist")
	}

	itemsInList, _ := shoppinglist.GetItems()

	Equal(t, items.Title, itemsInList[0].Title)
	Equal(t, items.Bought, itemsInList[0].Bought)
	Equal(t, items.ParentListID, itemsInList[0].ParentListID)
	Equal(t, items.Position, itemsInList[0].Position)
	Equal(t, items.ItemID, itemsInList[0].ItemID)
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
		items := &models.Item{
			ParentListID: id,
			ItemID:       itemID,
			Title:        util.StringWithCharset(100),
			Position:     position,
		}
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Items: *items,
			Owner: owner,
		}

		created, err := shoppinglist.Create(0, false)
		if err != nil || !created {
			t.Errorf("Failed to create shoppinglist %s", err.Error())
		}

		_, err = shoppinglist.AddItem()
		if err != nil {
			t.Errorf("Error while adding item: %s", err)
		}

		position2 := util.RandomPosition()
		items2 := &models.Item{
			ParentListID: id,
			ItemID:       itemID,
			Title:        util.StringWithCharset(100),
			Position:     position2,
		}
		shoppinglist2 := Shoppinglist{Items: *items2}

		_, err = shoppinglist2.AddItem()
		if err != nil {
			t.Errorf("Error while adding item: %s", err)
		}

		lastPosition, err := shoppinglist.GetLastPosition()
		if err != nil {
			t.Errorf("Error while getting last position: %s", err)
		}
		t.Logf("Position: %d", position)
		t.Logf("Position 2: %d", position2)
		t.Logf("Last Position: %d", lastPosition)

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

		created, err := shoppinglist.Create(0, false)
		if err != nil || !created {
			t.Errorf("Failed to create shoppinglist %s", err.Error())
		}

		lastPosition, _ := shoppinglist.GetLastPosition()

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
		items := &models.Item{
			ParentListID: id,
			ItemID:       itemID,
			Title:        itemTitle,
			Position:     position,
		}
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Items: *items,
			Owner: owner,
		}

		created, err := shoppinglist.Create(0, false)
		if err != nil {
			t.Errorf("Failed to create shoppinglist %s", err.Error())
		}

		_, err = shoppinglist.AddItem()
		if err != nil {
			t.Errorf("Error while adding item: %s", err)
		}

		item := Item{
			ItemID:       itemID,
			ParentListID: id,
		}

		i, err := item.GetItem()
		if err != nil {
			t.Errorf("Error while getting item: %s", err)
		}

		Equal(t, true, created)
		Equal(t, id, i.ParentListID)
		Equal(t, itemID, i.ItemID)
		Equal(t, itemTitle, i.Title)
		Equal(t, position, i.Position)
	})

	t.Run("Get Item that doesn't exist", func(t *testing.T) {
		id := util.RandomIntWithLength(500)
		itemID := util.RandomIntWithLength(200)

		item := Item{
			ItemID:       itemID,
			ParentListID: id,
		}

		i, err := item.GetItem()
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
	items := &models.Item{
		ParentListID: id,
		ItemID:       itemID,
		Title:        itemTitle,
		Position:     position,
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Items: *items,
		Owner: owner,
	}

	created, err := shoppinglist.Create(0, false)
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	_, err = shoppinglist.AddItem()
	if err != nil {
		t.Errorf("Error while adding item: %s", err)
	}

	item := Item{
		ItemID:       itemID,
		ParentListID: id,
	}

	i, err := item.GetItem()
	if err != nil {
		t.Errorf("Error while getting item: %s", err)
	}

	newTitle := util.StringWithCharset(10000)
	newPosition := util.RandomPosition()
	newBought := util.RandomBoolean()
	updatedItem := Item{
		ItemID:       itemID,
		ParentListID: id,
		Title:        &newTitle,
		Position:     &newPosition,
		Bought:       &newBought,
	}

	err = updatedItem.UpdateItem()
	if err != nil {
		t.Errorf("Error while updating item: %s", err)
	}

	i2, err := item.GetItem()
	if err != nil {
		t.Errorf("Error while getting item: %s", err)
	}

	Equal(t, true, created)
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

		created, err := shoppinglist.Create(0, false)
		if err != nil {
			t.Errorf("Failed to create shoppinglist %s", err.Error())
		}

		_, err = participant.AddParticipant()
		if err != nil {
			t.Errorf("Error while adding participant: %s", err)
		}

		Equal(t, true, created)
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

		_, err := participant.AddParticipant()

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

	created, err := shoppinglist.Create(0, false)
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	_, err = participant.AddParticipant()
	if err != nil {
		t.Errorf("Error while adding participant: %s", err)
	}

	participants, err := participant.GetParticipants()
	if err != nil {
		t.Errorf("Error while getting participants: %s", err)
	}

	Equal(t, true, created)
	Equal(t, participantEmail, participants[0].Email)
	Equal(t, participantStatus, participants[0].Status)
	Equal(t, id, participants[0].ParentListID)
}

func TestRemoveParticipant(t *testing.T) {
	Setup()

	t.Run("Remove Participant", func(t *testing.T) {
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

		created, err := shoppinglist.Create(0, false)
		if err != nil {
			t.Errorf("Failed to create shoppinglist %s", err)
		}

		_, err = participant.AddParticipant()
		if err != nil {
			t.Errorf("Error while adding participant: %s", err)
		}

		participants, err := participant.GetParticipants()
		if err != nil {
			t.Errorf("Error while getting participants: %s", err)
		}

		delParticipant := Participant{
			ID:           participants[0].ID,
			ParentListID: id,
			Email:        participantEmail,
			Status:       participantStatus,
		}

		delErr := delParticipant.RemoveParticipant()
		if delErr != nil {
			t.Errorf("Error while removing participant: %s", err)
		}

		participantsAfter, err := participant.GetParticipants()
		if err != nil {
			t.Errorf("Error while getting participants: %s", err)
		}

		if len(participantsAfter) > 0 {
			t.Errorf("Participant did not get deleted")
		}

		Equal(t, true, created)
	})

	t.Run("Remove Participant when Shoppinglist doesn't exist", func(t *testing.T) {
		id := util.RandomIntWithLength(50000)
		participantEmail := util.RandomEmail()
		participantStatus := util.StringWithCharset(20000)
		participant := Participant{
			ParentListID: id,
			Email:        participantEmail,
			Status:       participantStatus,
		}

		err := participant.RemoveParticipant()

		NotNil(t, err)
	})
}

func TestDeleteItem(t *testing.T) {
	Setup()

	id := util.RandomIntWithLength(50000)
	itemID := util.RandomIntWithLength(50000)
	title := "title3332999" + util.StringWithCharset(20000)
	owner := "owner999" + util.StringWithCharset(30000)
	items := models.Item{
		ParentListID: id,
		ItemID:       itemID,
		Title:        util.StringWithCharset(10000),
		Position:     util.RandomPosition(),
	}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		Items: items,
		Owner: owner,
	}

	created, err := shoppinglist.Create(0, false)
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	item, err := shoppinglist.AddItem()
	if err != nil {
		t.Errorf("Failed to edit shoppinglist %s", err.Error())
	}

	i := Item{
		ItemID:       itemID,
		ParentListID: id,
	}

	err = i.DeleteItem()
	if err != nil {
		t.Errorf("Error while deleting item: %s", err)
	}

	it, err := shoppinglist.GetItems()
	if err != nil {
		t.Errorf("Error while getting items: %s", err)
	}
	t.Log(it)

	if len(it) > 0 {
		t.Errorf("Item didn't get deleted")
	}

	Equal(t, true, created)
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
	items := []models.Item{
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

	created, err := shoppinglist.Create(0, false)
	if err != nil {
		t.Errorf("Error while creating shoppinglist: %s", err)
	}

	for _, item := range items {
		_, err = models.AddItem(item)
		if err != nil {
			t.Errorf("Error while adding item: %s", err)
		}
	}

	title3 := util.StringWithCharset(100)
	position3 := util.RandomPosition()
	title4 := util.StringWithCharset(100)
	position4 := util.RandomPosition()
	updateItems := []models.Item{
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

	i, err := shoppinglist.GetItems()
	if err != nil {
		t.Errorf("Error while getting items: %s", err)
	}

	Equal(t, true, created)
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

		created, err := shoppinglist.Create(0, false)
		if err != nil {
			t.Errorf("Error while creating shoppinglist: %s", err)
		}

		_, err = participant.AddParticipant()

		if err != nil {
			t.Errorf("Error while adding participant to list: %s", err)
		}

		lists, err := GetListsByParticipant(participantEmail)
		if err != nil {
			t.Errorf("Error while getting lists by participant: %s", err)
		}

		True(t, created)
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

		created, err := shoppinglist.Create(0, false)
		if err != nil {
			t.Errorf("Error while creating shoppinglist: %s", err)
		}

		_, err = participant.AddParticipant()
		if err != nil {
			t.Errorf("Error while adding participant to list: %s", err)
		}

		_, err = participant2.AddParticipant()
		if err != nil {
			t.Errorf("Error while adding participant to list: %s", err)
		}

		_, err = participant3.AddParticipant()
		if err != nil {
			t.Errorf("Error while adding participant to list: %s", err)
		}

		lists, err := GetListsByParticipant(participantEmail)
		if err != nil {
			t.Errorf("Error while getting lists by participant: %s", err)
		}

		True(t, created)
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
