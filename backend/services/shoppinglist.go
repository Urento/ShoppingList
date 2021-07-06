package services

import (
	"log"

	"github.com/urento/shoppinglist/models"
)

type Shoppinglist struct {
	ID           int
	Title        string
	Items        []string
	Owner        string
	Participants []string
	PageNum      int
	PageSize     int
}

func (s *Shoppinglist) Create() (created bool, err error) {
	shoppinglist := map[string]interface{}{
		"id":           s.ID,
		"title":        s.Title,
		"items":        s.Items,
		"owner":        s.Owner,
		"participants": s.Participants,
	}

	if err := models.CreateList(shoppinglist); err != nil {
		log.Fatal(err.Error())
		return false, err
	}

	return true, nil
}

func (s *Shoppinglist) Edit() error {
	return models.EditList(s.ID, map[string]interface{}{
		"id":           s.ID,
		"title":        s.Title,
		"items":        s.Items,
		"owner":        s.Owner,
		"participants": s.Participants,
	})
}

func (s *Shoppinglist) GetList() (*models.Shoppinglist, error) {
	shoppinglist, err := models.GetList(s.ID)
	if err != nil {
		return nil, err
	}
	return shoppinglist, nil
}

func (s *Shoppinglist) Delete() error {
	return models.DeleteList(s.ID)
}

func (s *Shoppinglist) ExistsByID() (bool, error) {
	return models.ExistByID(s.ID)
}

func DropTable() {
	models.DropTable()
}
