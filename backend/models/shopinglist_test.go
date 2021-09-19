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

	t.Run("TestGetTotalListsByOwner", func(t *testing.T) {
		id := RandomInt()
		title := "title" + StringWithCharset(20)
		owner := "ParentListID" + StringWithCharset(30)
		participants := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
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
		id := RandomInt()
		title := "title" + StringWithCharset(20)
		owner := "ParentListID" + StringWithCharset(30)
		participants := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
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

		id2 := RandomInt()
		title2 := "title" + StringWithCharset(20)
		participants2 := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
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

		id3 := RandomInt()
		title3 := "title" + StringWithCharset(20)
		participants3 := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
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

	id := RandomInt()
	title := "title" + StringWithCharset(20)
	owner := "Owner123123123123" + StringWithCharset(30)
	participants := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	shoppinglist := map[string]interface{}{
		"id":           id,
		"title":        title,
		"owner":        owner,
		"participants": participants,
	}

	if err := CreateList(shoppinglist); err != nil {
		t.Errorf("Error while creating Shoppinglist 1 %s", err.Error())
	}

	id2 := RandomInt()
	title2 := "title" + StringWithCharset(20)
	participants2 := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
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

func StringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomInt() int {
	return seededRand.Intn(90000)
}

func SetupTest() {
	Setup()
	util.Setup()
}
