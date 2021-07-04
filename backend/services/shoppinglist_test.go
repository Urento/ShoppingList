package services

import (
	"testing"
)

func TestCreate(t *testing.T) {
	id := 3498757245897
	title := "Test Title öikuhadsfbglhkjdsbnfgjlhdsfbgjhdksfbg"
	items := []string{"1", "2", "3", "skjdhfgbnsdljhkfgbdkhjsfbgjdhfbgjdhfg"}
	owner := "Urento hkalsbdrgojuzsdafghkdsfg"
	participants := []string{"Me", "you", "m123123123e", "38383883"}
	shoppinglist := Shoppinglist{
		ID:           id,
		Title:        title,
		Items:        items,
		Owner:        owner,
		Participants: participants,
	}
	err := shoppinglist.Create()
	if err != nil {
		t.Errorf("Failed to create shoppinglist")
	}

	list := Shoppinglist{ID: id}
	exists, err := list.ExistsByID()
	if err != nil || !exists {
		t.Errorf("Shoppinglist did not get created")
	}

	l, err := list.GetList()
	if err != nil {
		t.Errorf("Shoppinglist not found")
	}

	if l.ID != id || !testEq(l.Items, items) || l.Owner != owner || !testEq(l.Participants, participants) || l.Title != title {
		t.Errorf("Values weren't correct ")
	}

}

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
