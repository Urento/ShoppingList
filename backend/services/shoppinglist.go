package services

import (
	"fmt"
	"time"

	"github.com/urento/shoppinglist/models"
)

type Shoppinglist struct {
	ID           int
	Title        string
	Items        models.Item
	Owner        string
	Participants []*models.Participant
}

type Item struct {
	ItemID       int
	Title        *string
	Position     *int64
	Bought       *bool
	ParentListID int
}

type Participant struct {
	ID           int
	ParentListID int
	Status       string
	Email        string
}

func (s *Shoppinglist) Create(userId int, withNotification bool) (bool, error) {
	shoppinglist := models.Shoppinglist{
		ID:    s.ID,
		Title: s.Title,
		Owner: s.Owner,
	}

	if err := models.CreateList(shoppinglist); err != nil {
		return false, err
	}

	if withNotification {
		notification := models.Notification{
			UserID:           userId,
			Title:            "New Shoppinglist",
			Text:             fmt.Sprintf("%s was created", s.Title),
			NotificationType: "new_shoppinglist",
			Date:             time.Now().Format("02.01.2006"),
		}

		if err := models.CreateNotification(notification); err != nil {
			return false, err
		}
	}

	return true, nil
}

func (s *Shoppinglist) Edit(userId int, withNotification bool) error {
	shoppinglist := map[string]interface{}{
		"id":    s.ID,
		"title": s.Title,
		"items": s.Items,
		"owner": s.Owner,
	}

	if withNotification {
		notification := models.Notification{
			UserID:           userId,
			Title:            "New Shoppinglist",
			Text:             fmt.Sprintf("%s was created", s.Title),
			NotificationType: "new_shoppinglist",
			Date:             time.Now().Format("02.01.2006"),
		}

		if err := models.CreateNotification(notification); err != nil {
			return err
		}
	}

	return models.EditList(s.ID, shoppinglist)
}

func (s *Shoppinglist) GetList() (*models.Shoppinglist, error) {
	return models.GetList(s.ID, s.Owner)
}

func (s *Shoppinglist) GetListsByOwner() (*[]models.Shoppinglist, error) {
	return models.GetListByEmail(s.Owner)
}

func (s *Shoppinglist) GetLastPosition() (int64, error) {
	return models.GetLastPosition(s.ID)
}

func (s *Shoppinglist) Delete(userId int, withNotification bool) error {
	if withNotification {
		notification := models.Notification{
			UserID:           userId,
			Title:            "Deleted Shoppinglist",
			Text:             "Shoppinglist was deleted", //add name
			NotificationType: "deleted_shoppinglist",
			Date:             time.Now().Format("02.01.2006"),
		}

		if err := models.CreateNotification(notification); err != nil {
			return err
		}
	}

	return models.DeleteList(s.ID)
}

func (s *Shoppinglist) ExistsByID() (bool, error) {
	return models.ExistByID(s.ID)
}

func (s *Shoppinglist) GetItems() ([]models.Item, error) {
	return models.GetItems(s.ID)
}

func (i *Item) GetItem() (models.Item, error) {
	return models.GetItem(i.ParentListID, i.ItemID)
}

func (i *Item) UpdateItem() error {
	item := models.Item{
		ItemID:       i.ItemID,
		ParentListID: i.ParentListID,
		Title:        string(*i.Title),
		Position:     int64(*i.Position),
		Bought:       bool(*i.Bought),
	}
	return models.UpdateItem(item)
}

func (s *Shoppinglist) AddItem() (*models.Item, error) {
	return models.AddItem(s.Items)
}

func (s *Item) DeleteItem() error {
	return models.DeleteItem(s.ParentListID, s.ItemID)
}

func (p *Participant) AddParticipant() (models.Participant, error) {
	participant := models.Participant{
		ParentListID: p.ParentListID,
		Status:       p.Status,
		Email:        p.Email,
	}
	return models.AddParticipant(participant)
}

func (p *Participant) RemoveParticipant() error {
	return models.RemoveParticipant(p.ParentListID, p.ID)
}

func (p *Participant) GetParticipants() ([]models.Participant, error) {
	return models.GetParticipants(p.ParentListID)
}
