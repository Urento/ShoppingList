package services

import "github.com/urento/shoppinglist/models"

type Shoppinglist struct {
	ID           int
	ListID       int
	Title        string
	Items        []string
	Owner        string
	Participants []string
	PageNum      int
	PageSize     int
}

func (s *Shoppinglist) Create() error {
	shoppinglist := map[string]interface{}{
		"list_id":      s.ListID,
		"title":        s.Title,
		"items":        s.Items,
		"owner":        s.Owner,
		"participants": s.Participants,
	}

	if err := models.CreateList(shoppinglist); err != nil {
		return err
	}

	return nil
}

func (s *Shoppinglist) Edit() error {
	return models.EditList(s.ID, map[string]interface{}{
		"list_id":      s.ListID,
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

func (s *Shoppinglist) Exists() (bool, error) {
	return models.ExistByID(s.ID)
}
