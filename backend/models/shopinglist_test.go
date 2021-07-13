package models

import (
	"math/rand"
	"strings"
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
	t.Log(owner)

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
	} else {
		t.Log("no error")
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

	t.Log(owner)
	t.Log(lists)

	if len(lists) <= 0 {
		t.Errorf("List is Empty")
	}

	containsOwner := strings.Contains(lists[0].Owner, owner)
	containsOwner1 := strings.Contains(lists[1].Owner, owner)

	Equal(t, containsOwner, true)
	Equal(t, containsOwner1, true)
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
