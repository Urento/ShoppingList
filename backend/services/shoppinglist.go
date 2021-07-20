package services

import (
	"github.com/urento/shoppinglist/models"
)

type Shoppinglist struct {
	ID           int
	Title        string
	Items        []string
	Owner        string
	Participants []string
	Position     int //TODO
	PageNum      int
	PageSize     int
}

//TODO: Encrypt Emails
//TODO: Add Order to Items

func (s *Shoppinglist) Create() (created bool, err error) {
	shoppinglist := map[string]interface{}{
		"id":           s.ID,
		"title":        s.Title,
		"items":        s.Items,
		"owner":        s.Owner,
		"position":     s.Position,
		"participants": s.Participants,
	}

	if err := models.CreateList(shoppinglist); err != nil {
		return false, err
	}

	return true, nil
}

func (s *Shoppinglist) Edit() error {
	shoppinglist := map[string]interface{}{
		"title":        s.Title,
		"items":        s.Items,
		"owner":        s.Owner,
		"position":     s.Position,
		"participants": s.Participants,
	}
	return models.EditList(s.ID, shoppinglist)
}

func (s *Shoppinglist) GetList() (*models.Shoppinglist, error) {
	shoppinglist, err := models.GetList(s.ID)
	if err != nil {
		return nil, err
	}
	return shoppinglist, nil
}

func (s *Shoppinglist) GetLastPosition() (int64, error) {
	//TODO
	return 0, nil
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
