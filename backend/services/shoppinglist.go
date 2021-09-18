package services

import (
	"log"

	"github.com/urento/shoppinglist/models"
)

type Shoppinglist struct {
	ID           int
	Title        string
	Items        models.Item
	Owner        string
	Participants []string
	PageNum      int
	PageSize     int
}

func (s *Shoppinglist) Create() (created bool, err error) {
	shoppinglist := map[string]interface{}{
		"id":           s.ID,
		"title":        s.Title,
		"owner":        s.Owner,
		"participants": s.Participants,
	}

	if err := models.CreateList(shoppinglist); err != nil {
		return false, err
	}

	return true, nil
}

func (s *Shoppinglist) Edit() error {
	shoppinglist := map[string]interface{}{
		"id":           s.ID,
		"title":        s.Title,
		"items":        s.Items,
		"owner":        s.Owner,
		"participants": s.Participants,
	}

	return models.EditList(s.ID, shoppinglist)
}

func (s *Shoppinglist) GetList() (*models.Shoppinglist, error) {
	shoppinglist, err := models.GetList(s.ID, s.Owner)
	return shoppinglist, err
}

func (s *Shoppinglist) GetListsByOwner() (*[]models.Shoppinglist, error) {
	shoppinglist, err := models.GetListByEmail(s.Owner)
	return shoppinglist, err
}

func (s *Shoppinglist) SendInvitationEmails() error {
	for _, val := range s.Participants {
		log.Print(val)
		//TODO: Send Emails
	}
	return nil
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

func (s *Shoppinglist) GetItems() ([]models.Item, error) {
	return models.GetItems(s.ID, s.Owner)
}

func (s *Shoppinglist) AddItem() (*models.Item, error) {
	return models.AddItem(s.Items)
}
