package services

import (
	"fmt"
	"math/rand"
	"strings"
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

	id := randInt(1000, 50000)
	title := "title" + StringWithCharset(20)
	items := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	owner := "owner" + StringWithCharset(30)
	participants := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	shoppinglist := Shoppinglist{
		ID:           id,
		Title:        title,
		Items:        items,
		Owner:        owner,
		Participants: participants,
	}
	created, err := shoppinglist.Create()
	if err != nil || !created {
		fmt.Println(err.Error())
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

	Equal(t, created, true)
	Equal(t, id, l.ID)
	Equal(t, items, l.Items)
	Equal(t, participants, l.Participants)
	Equal(t, title, l.Title)
	Equal(t, owner, l.Owner)
}

func TestExistsByID(t *testing.T) {
	Setup()

	id := randInt(1000, 50000)
	title := "title" + StringWithCharset(20)
	items := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	owner := "owner" + StringWithCharset(30)
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

	list := Shoppinglist{ID: id}
	exists, err := list.ExistsByID()
	if err != nil || !exists {
		t.Errorf("Shoppinglist did not get created %s", err.Error())
	}

	Equal(t, created, true)
	Equal(t, exists, true)
}

func TestCreateAndDelete(t *testing.T) {
	Setup()

	id := randInt(1000, 50000)
	title := "title" + StringWithCharset(20)
	items := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	owner := "owner" + StringWithCharset(30)
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

	list := Shoppinglist{ID: id}
	exists, err := list.ExistsByID()
	if err != nil || !exists {
		t.Errorf("Shoppinglist did not get created %s", err.Error())
	}

	l, err := list.GetList()
	if err != nil {
		t.Errorf("Shoppinglist not found %s", err.Error())
	}

	delErr := list.Delete()
	if delErr != nil {
		t.Errorf("Shoppinglist couldn't be deleted %s", err.Error())
	}

	Equal(t, created, true)
	Equal(t, exists, true)
	Equal(t, id, l.ID)
	Equal(t, items, l.Items)
	Equal(t, participants, l.Participants)
	Equal(t, title, l.Title)
	Equal(t, owner, l.Owner)
	Equal(t, delErr, nil)
}

func TestCreateAndEdit(t *testing.T) {
	Setup()

	id := randInt(1000, 50000)
	title := "title" + StringWithCharset(20)
	items := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	owner := "owner" + StringWithCharset(30)
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

	title2 := "title2" + StringWithCharset(20)
	items2 := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	owner2 := "owner2" + StringWithCharset(30)
	participants2 := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	shoppinglist = Shoppinglist{
		ID:           id,
		Title:        title2,
		Items:        items2,
		Owner:        owner2,
		Participants: participants2,
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

	Equal(t, created, true)
	Equal(t, id, l.ID)
	Equal(t, items2, l.Items)
	Equal(t, participants2, l.Participants)
	Equal(t, title2, l.Title)
	Equal(t, owner2, l.Owner)
	NotEqual(t, items, l.Items)
	NotEqual(t, owner, l.Owner)
	NotEqual(t, title, l.Title)
	NotEqual(t, participants, l.Participants)
}

func TestDuplicationError(t *testing.T) {
	Setup()
	id1 := randInt(1000, 50000)
	shoppinglist1 := ShoppinglistObject(id1)

	created, err := shoppinglist1.Create()
	if err != nil || !created {
		t.Errorf("Error while creating Shoppinglist %s", err.Error())
	}

	_, err = shoppinglist1.Create()

	containsKeyConstraint := strings.Contains(err.Error(), "violates unique constraint")

	NotEqual(t, err, nil)
	Equal(t, containsKeyConstraint, true)
}

//TODO: Use this to avoid code duplication
func ShoppinglistObject(id int) Shoppinglist {
	//id := randInt(1000, 50000)
	title := "title" + StringWithCharset(20)
	items := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	owner := "owner" + StringWithCharset(30)
	participants := pq.StringArray{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	shoppinglist := Shoppinglist{
		ID:           id,
		Title:        title,
		Items:        items,
		Owner:        owner,
		Participants: participants,
	}
	return shoppinglist
}

func StringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(100000)
}

func Setup() {
	models.Setup()
	util.Setup()
}
