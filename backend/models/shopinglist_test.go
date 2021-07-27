package models

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/pkg/util"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func TestGetTotalListsByOwner(t *testing.T) {
	SetupTest()

	id := seededRand.Intn(90000)
	title := "title" + StringWithCharset(20)
	items := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	owner := "owner" + StringWithCharset(30)
	participants := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	position := seededRand.Intn(5000)
	shoppinglist := map[string]interface{}{
		"id":           id,
		"title":        title,
		"items":        items,
		"owner":        owner,
		"position":     position,
		"participants": participants,
	}

	err := CreateList(shoppinglist)
	if err != nil {
		t.Errorf("Error while creating Shoppinglist %s", err.Error())
	}

	count, err := GetTotalListsByOwner(owner)
	if err != nil {
		t.Errorf("Error while getting the total lists by owner %s", err.Error())
	}

	Equal(t, count, int64(1))
}

func TestGetTotalListsByOwnerWithMultipleLists(t *testing.T) {
	SetupTest()

	id := seededRand.Intn(90000)
	title := "title" + StringWithCharset(20)
	items := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	owner := "owner" + StringWithCharset(30)
	participants := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	position := seededRand.Intn(5000)
	shoppinglist := map[string]interface{}{
		"id":           id,
		"title":        title,
		"items":        items,
		"owner":        owner,
		"position":     position,
		"participants": participants,
	}

	err := CreateList(shoppinglist)
	if err != nil {
		t.Errorf("Error while creating Shoppinglist 1 %s", err.Error())
	}

	id2 := seededRand.Intn(90000)
	title2 := "title" + StringWithCharset(20)
	items2 := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	participants2 := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	position2 := seededRand.Intn(5000)
	shoppinglist2 := map[string]interface{}{
		"id":           id2,
		"title":        title2,
		"items":        items2,
		"owner":        owner,
		"position":     position2,
		"participants": participants2,
	}

	err = CreateList(shoppinglist2)
	if err != nil {
		t.Errorf("Error while creating Shoppinglist 2 %s", err.Error())
	}

	id3 := seededRand.Intn(90000)
	title3 := "title" + StringWithCharset(20)
	items3 := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	participants3 := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	position3 := seededRand.Intn(5000)
	shoppinglist3 := map[string]interface{}{
		"id":           id3,
		"title":        title3,
		"items":        items3,
		"owner":        owner,
		"position":     position3,
		"participants": participants3,
	}

	err = CreateList(shoppinglist3)
	if err != nil {
		t.Errorf("Error while creating Shoppinglist 3 %s", err.Error())
	}

	count, err := GetTotalListsByOwner(owner)
	if err != nil {
		t.Errorf("Error while getting the total lists by owner %s", err.Error())
	}

	Equal(t, count, int64(3))
}

func TestGetListsByOwner(t *testing.T) {
	SetupTest()

	id := seededRand.Intn(90000)
	title := "title" + StringWithCharset(20)
	items := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	owner := "owner123123123123" + StringWithCharset(30)
	participants := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	position := seededRand.Intn(100)
	shoppinglist := map[string]interface{}{
		"id":           id,
		"title":        title,
		"items":        items,
		"owner":        owner,
		"position":     position,
		"participants": participants,
	}

	if err := CreateList(shoppinglist); err != nil {
		t.Errorf("Error while creating Shoppinglist 1 %s", err.Error())
	}

	id2 := seededRand.Intn(10000)
	title2 := "title" + StringWithCharset(20)
	items2 := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	participants2 := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	position2 := seededRand.Intn(100)
	shoppinglist2 := map[string]interface{}{
		"id":           id2,
		"title":        title2,
		"items":        items2,
		"owner":        owner,
		"position":     position2,
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
	Equal(t, position, lists[0].Position)
	Equal(t, id, lists[0].ID)
	Equal(t, title2, lists[1].Title)
	Equal(t, position2, lists[1].Position)
	Equal(t, id2, lists[1].ID)
}

func TestExistsByID(t *testing.T) {
	SetupTest()

	id := seededRand.Intn(90000)
	title := "title" + StringWithCharset(20)
	items := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	owner := "owner123123123123" + StringWithCharset(30)
	participants := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	position := seededRand.Intn(100)
	shoppinglist := map[string]interface{}{
		"id":           id,
		"title":        title,
		"items":        items,
		"owner":        owner,
		"position":     position,
		"participants": participants,
	}

	if err := CreateList(shoppinglist); err != nil {
		t.Errorf("Error while creating Shoppinglist 1 %s", err.Error())
	}

	exists, err := ExistByID(id)
	if err != nil {
		t.Errorf("Error while checking if the shoppinglist exists by id %s", err)
	}

	Equal(t, true, exists)
	Equal(t, nil, err)
}

func StringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func SetupTest() {
	Setup()
	util.Setup()
}
