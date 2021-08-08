package services

import (
	"math/rand"
	"testing"
	"time"

	"github.com/lib/pq"
	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/models"
	"github.com/urento/shoppinglist/pkg/util"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func TestCreateAndCheck(t *testing.T) {
	Setup()

	id := seededRand.Intn(90000)
	title := "title" + StringWithCharset(20)
	/*items := []models.Item{
		{
			ParentListID: RandomInt(),
			Title:        StringWithCharset(10),
			Position:     RandomInt(),
		},
		{
			ParentListID: RandomInt(),
			Title:        StringWithCharset(10),
			Position:     RandomInt(),
		},
		{
			ParentListID: RandomInt(),
			Title:        StringWithCharset(10),
			Position:     RandomInt(),
		},
		{
			ParentListID: RandomInt(),
			Title:        StringWithCharset(10),
			Position:     RandomInt(),
		},
	}*/
	owner := "owner" + StringWithCharset(30)
	participants := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		//Items:        items,
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
	//Equal(t, items, l.Items)
	Equal(t, participants, l.Participants)
	Equal(t, title, l.Title)
	Equal(t, owner, l.Owner)
	Equal(t, nil, err)
	Equal(t, true, created)
}

func TestExistsByID(t *testing.T) {
	Setup()

	id := seededRand.Intn(90000)
	title := "titlesdfgdsghdshgfdzhjf" + StringWithCharset(20)
	/*items := []models.Item{
		{
			ParentListID: RandomInt(),
			Title:        StringWithCharset(10),
			Position:     RandomInt(),
		},
		{
			ParentListID: RandomInt(),
			Title:        StringWithCharset(10),
			Position:     RandomInt(),
		},
		{
			ParentListID: RandomInt(),
			Title:        StringWithCharset(10),
			Position:     RandomInt(),
		},
		{
			ParentListID: RandomInt(),
			Title:        StringWithCharset(10),
			Position:     RandomInt(),
		},
	}*/
	owner := "ownersthfdghdfhfdthfxgdh" + StringWithCharset(30)
	participants := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		//Items:        items,
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

func TestCreateAndEdit(t *testing.T) {
	Setup()

	id := RandomInt()
	title := "title3332999" + StringWithCharset(20)
	owner := "owner999" + StringWithCharset(30)
	participants := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
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

	title2 := "title2111111999" + StringWithCharset(20)
	/*items := []models.Item{
		{
			ParentListID: RandomInt(),
			Title:        StringWithCharset(10),
			Position:     RandomInt(),
		},
		{
			ParentListID: RandomInt(),
			Title:        StringWithCharset(10),
			Position:     RandomInt(),
		},
		{
			ParentListID: RandomInt(),
			Title:        StringWithCharset(10),
			Position:     RandomInt(),
		},
		{
			ParentListID: RandomInt(),
			Title:        StringWithCharset(10),
			Position:     RandomInt(),
		},
	}*/
	participants2 := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	shoppinglist = Shoppinglist{
		ID:    id,
		Title: title2,
		//	Items:        items,
		Owner:        owner,
		Participants: participants2,
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
	//Equal(t, items, l.Items)
	Equal(t, participants2, l.Participants)
	Equal(t, title2, l.Title)
	Equal(t, owner, l.Owner)
	NotEqual(t, title, l.Title)
	NotEqual(t, participants, l.Participants)
	Equal(t, nil, err)
}

func TestAddItem(t *testing.T) {
	Setup()

	id := RandomInt()
	title := "title3332999" + StringWithCharset(20)
	owner := "owner999" + StringWithCharset(30)
	items := models.Item{
		ParentListID: id,
		Title:        StringWithCharset(10),
		Position:     RandomInt(),
	}
	participants := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	shoppinglist := Shoppinglist{
		ID:           id,
		Title:        title,
		Items:        items,
		Owner:        owner,
		Participants: participants,
	}

	created, err := shoppinglist.Create()
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	err = shoppinglist.AddItem()
	if err != nil {
		t.Errorf("Failed to edit shoppinglist %s", err.Error())
	}

	Equal(t, false, created)
	Equal(t, nil, err)
}

func TestGetItems(t *testing.T) {
	Setup()

	id := seededRand.Intn(90000)
	title := "title3332999" + StringWithCharset(20)
	/*items := []models.Item{
		{
			Title:    StringWithCharset(10),
			Position: RandomInt(),
		},
		{
			Title:    StringWithCharset(10),
			Position: RandomInt(),
		},
		{
			Title:    StringWithCharset(10),
			Position: RandomInt(),
		},
		{
			Title:    StringWithCharset(10),
			Position: RandomInt(),
		},
	}*/
	owner := "owner999" + StringWithCharset(30)
	participants := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	shoppinglist := Shoppinglist{
		ID:    id,
		Title: title,
		//Items:        items,
		Owner:        owner,
		Participants: participants,
	}

	created, err := shoppinglist.Create()
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	if !created {
		t.Errorf("Error while creating shoppinglist")
	}

	/*itemsInList, err := shoppinglist.GetItems()
	if err != nil {
		t.Errorf("Error while getting items: %s", err)
	}*/

	Equal(t, nil, err)
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

func Setup() {
	models.Setup()
	util.Setup()
}
