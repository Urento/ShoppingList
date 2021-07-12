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
	items := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	owner := "owner" + StringWithCharset(30)
	participants := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	position := seededRand.Intn(5000)
	shoppinglist := Shoppinglist{
		ID:           id,
		Title:        title,
		Items:        items,
		Owner:        owner,
		Participants: participants,
		Position:     position,
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

	Equal(t, created, true)
	Equal(t, id, l.ID)
	Equal(t, items, l.Items)
	Equal(t, participants, l.Participants)
	Equal(t, title, l.Title)
	Equal(t, owner, l.Owner)
	Equal(t, position, l.Position)
	Equal(t, err, nil)
}

func TestExistsByID(t *testing.T) {
	Setup()

	id := seededRand.Intn(90000)
	title := "titlesdfgdsghdshgfdzhjf" + StringWithCharset(20)
	items := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	owner := "ownersthfdghdfhfdthfxgdh" + StringWithCharset(30)
	participants := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	position := seededRand.Intn(5000)
	shoppinglist := Shoppinglist{
		ID:           id,
		Title:        title,
		Items:        items,
		Owner:        owner,
		Participants: participants,
		Position:     position,
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

	Equal(t, created, true)
	Equal(t, exists, true)
	Equal(t, err, nil)
}

func TestCreateAndEdit(t *testing.T) {
	Setup()

	id := seededRand.Intn(90000)
	title := "title3332999" + StringWithCharset(20)
	items := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	owner := "owner999" + StringWithCharset(30)
	participants := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	position := seededRand.Intn(5000)
	shoppinglist := Shoppinglist{
		ID:           id,
		Title:        title,
		Items:        items,
		Owner:        owner,
		Participants: participants,
		Position:     position,
	}

	created, err := shoppinglist.Create()
	if err != nil {
		t.Errorf("Failed to create shoppinglist %s", err.Error())
	}

	title2 := "title2111111999" + StringWithCharset(20)
	items2 := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	owner2 := "owner2123123" + StringWithCharset(30)
	participants2 := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	position2 := seededRand.Intn(5000)
	shoppinglist = Shoppinglist{
		ID:           id,
		Title:        title2,
		Items:        items2,
		Owner:        owner2,
		Participants: participants2,
		Position:     position2,
	}

	err = shoppinglist.Edit()
	if err != nil {
		t.Errorf("Failed to edit shoppinglist %s", err.Error())
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

	Equal(t, created, true)
	Equal(t, id, l.ID)
	Equal(t, items2, l.Items)
	Equal(t, participants2, l.Participants)
	Equal(t, title2, l.Title)
	Equal(t, owner2, l.Owner)
	Equal(t, position2, l.Position)
	NotEqual(t, items, l.Items)
	NotEqual(t, owner, l.Owner)
	NotEqual(t, title, l.Title)
	NotEqual(t, participants, l.Participants)
	NotEqual(t, position, l.Position)
	Equal(t, err, nil)
}

func StringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func Setup() {
	models.Setup()
	util.Setup()
}
