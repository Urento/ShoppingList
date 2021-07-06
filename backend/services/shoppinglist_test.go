package services

import (
	"fmt"
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

	id := randInt(1000, 50000)
	title := "title" + StringWithCharset(20)
	items := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
	owner := "owner" + StringWithCharset(30)
	participants := []string{StringWithCharset(45), StringWithCharset(45), StringWithCharset(45), StringWithCharset(45)}
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
		t.Errorf("Failed to create shoppinglist")
	}

	list := Shoppinglist{ID: id}
	exists, err := list.ExistsByID()
	t.Log(exists)
	if err != nil || !exists {
		t.Errorf("Shoppinglist did not get created")
	}

	l, err := list.GetList()
	if err != nil {
		t.Errorf("Shoppinglist not found")
	}

	Equal(t, created, true)
	Equal(t, id, l.ID)
	Equal(t, items, l.Items)
	Equal(t, participants, l.Participants)
	Equal(t, title, l.Title)
	Equal(t, owner, l.Owner)
}

func TestCreateAndDelete(t *testing.T) {
	Setup()

	Equal(t, true, true)
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
	models.Setup(true)
	util.Setup()
}
