package services

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/models"
	"github.com/urento/shoppinglist/pkg/util"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func TestCreate(t *testing.T) {
	Setup()

	t.Run("Create and Check", func(t *testing.T) {
		id := RandomInt()
		title := "title" + StringWithCharset(20)
		owner := "owner" + StringWithCharset(30)
		participants := []*models.Participant{
			{
				ParentListID: id,
				Email:        util.RandomEmail(),
			},
			{
				ParentListID: id,
				Email:        util.RandomEmail(),
			},
		}
		shoppinglist := Shoppinglist{
			ID:           id,
			Title:        title,
			Owner:        owner,
			Participants: participants,
		}
		created, err := shoppinglist.Create()
		if err != nil || !created {
			t.Errorf("Failed to create shoppinglist %s", err.Error())
		}

		list := Shoppinglist{ID: id}
		exists, err := list.ExistsByID()
		if err != nil || !exists {
			t.Errorf("Shoppinglist did not get created %s", err.Error())
		}

		l, err := list.GetList()
		if err != nil {
			t.Errorf("Shoppinglist not found %s", err.Error())
		}

		err = list.Delete()
		if err != nil {
			t.Errorf("Shoppinglist couldn't be deleted %s", err.Error())
		}

		Equal(t, id, l.ID)
		Equal(t, participants, l.Participants)
		Equal(t, title, l.Title)
		Equal(t, owner, l.Owner)
		Equal(t, nil, err)
		Equal(t, true, created)
	})

	t.Run("Create and Edit", func(t *testing.T) {
		id := RandomInt()
		title := "title3332999" + StringWithCharset(20)
		owner := "owner999" + StringWithCharset(30)
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Owner: owner,
		}

		created, err := shoppinglist.Create()
		if err != nil {
			t.Errorf("Failed to create shoppinglist %s", err.Error())
		}

		title2 := "title2111111999" + StringWithCharset(20)
		shoppinglist = Shoppinglist{
			ID:    id,
			Title: title2,
			Owner: owner,
		}

		err = shoppinglist.Edit()
		if err != nil {
			t.Errorf("Failed to edit shoppinglist %s", err.Error())
		}

		list := Shoppinglist{ID: id, Owner: owner}
		exists, err := list.ExistsByID()
		if err != nil || !exists {
			t.Errorf("Shoppinglist did not get created %s", err.Error())
		}

		l, err := list.GetList()
		if err != nil {
			t.Errorf("Shoppinglist not found %s", err.Error())
		}

		err = list.Delete()
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

	id := RandomInt()
	title := "titlesdfgdsghdshgfdzhjf" + StringWithCharset(20)
	owner := "ownersthfdghdfhfdthfxgdh" + StringWithCharset(30)
	participants := []*models.Participant{
		{
			ParentListID: id,
			Email:        util.RandomEmail(),
		},
		{
			ParentListID: id,
			Email:        util.RandomEmail(),
		},
	}
	shoppinglist := Shoppinglist{
		ID:           id,
		Title:        title,
		Owner:        owner,
		Participants: participants,
	}

	created, err := shoppinglist.Create()
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	list := Shoppinglist{ID: id}
	exists, err := list.ExistsByID()
	if err != nil || !exists {
		t.Errorf("Shoppinglist did not get created %s", err.Error())
	}

	err = list.Delete()
	if err != nil {
		t.Errorf("Shoppinglist couldn't be deleted %s", err.Error())
	}

	Equal(t, true, created)
	Equal(t, true, exists)
	Equal(t, nil, err)
}

func TestAddItem(t *testing.T) {
	Setup()

	id := RandomInt()
	itemID := RandomInt()
	title := "title3332999" + StringWithCharset(20)
	owner := "owner999" + StringWithCharset(30)
	items := &models.Item{
		ParentListID: id,
		ItemID:       itemID,
		Title:        StringWithCharset(10),
		Position:     RandomPosition(),
	}
	participants := []*models.Participant{
		{
			ParentListID: id,
			Email:        util.RandomEmail(),
		},
		{
			ParentListID: id,
			Email:        util.RandomEmail(),
		},
	}
	shoppinglist := Shoppinglist{
		ID:           id,
		Title:        title,
		Items:        *items,
		Owner:        owner,
		Participants: participants,
	}

	created, err := shoppinglist.Create()
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

func TestGetItems(t *testing.T) {
	Setup()

	id := RandomInt()
	itemID := RandomInt()
	title := "title3332999" + StringWithCharset(20)
	owner := "owner999" + StringWithCharset(30)
	items := &models.Item{
		ParentListID: id,
		ItemID:       itemID,
		Title:        StringWithCharset(10),
		Position:     RandomPosition(),
	}
	participants := []*models.Participant{
		{
			ParentListID: id,
			Email:        util.RandomEmail(),
		},
		{
			ParentListID: id,
			Email:        util.RandomEmail(),
		},
	}
	shoppinglist := Shoppinglist{
		ID:           id,
		Title:        title,
		Items:        *items,
		Owner:        owner,
		Participants: participants,
	}

	created, err := shoppinglist.Create()
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

	t.Run("TestGetLastPositionWithTwoItems", func(t *testing.T) {
		id := RandomInt()
		itemID := RandomInt()
		title := "title3332999" + StringWithCharset(20)
		owner := "owner999" + StringWithCharset(30)
		position := RandomPosition()
		items := &models.Item{
			ParentListID: id,
			ItemID:       itemID,
			Title:        StringWithCharset(10),
			Position:     position,
		}
		shoppinglist := Shoppinglist{
			ID:    id,
			Title: title,
			Items: *items,
			Owner: owner,
		}

		created, err := shoppinglist.Create()
		if err != nil {
			t.Errorf("Failed to create shoppinglist %s", err.Error())
		}

		_, err = shoppinglist.AddItem()
		if err != nil {
			t.Errorf("Error while adding item: %s", err)
		}

		items2 := &models.Item{
			ParentListID: id,
			ItemID:       itemID,
			Title:        StringWithCharset(10),
			Position:     position,
		}
		shoppinglist2 := Shoppinglist{
			ID:    id,
			Items: *items2,
			Owner: owner,
		}

		_, err = shoppinglist2.AddItem()
		if err != nil {
			t.Errorf("Error while adding item: %s", err)
		}

		if !created {
			t.Errorf("Error while creating shoppinglist")
		}

		lastPosition, err := shoppinglist.GetLastPosition()
		if err != nil {
			t.Errorf("Error while getting last position: %s", err)
		}
		t.Log(lastPosition)

		Equal(t, position, lastPosition)
	})

	t.Run("TestGetLastPositionWithMultipleItems", func(t *testing.T) {

	})

	t.Run("TestGetLastPositionWithoutAnyItems", func(t *testing.T) {

	})
}

func StringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomInt() int {
	s1 := rand.NewSource(time.Now().Unix())
	r1 := rand.New(s1)
	return r1.Intn(900000)
}

func RandomPosition() int {
	s1 := rand.NewSource(time.Now().Unix())
	r1 := rand.New(s1)
	return r1.Intn(10)
}

func Setup() {
	models.Setup()
	util.Setup()
}
