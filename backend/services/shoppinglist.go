package services

import (
	"fmt"
	"log"

	"github.com/urento/shoppinglist/models"
)

type Shoppinglist struct {
	ID           int
	Title        string
	Items        models.Item
	Owner        string
	Participants []*models.Participant
	PageNum      int
	PageSize     int
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
	shoppinglist := map[string]interface{}{
		"id":           s.ID,
		"title":        s.Title,
		"owner":        s.Owner,
		"participants": s.Participants,
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
		}

		if err := models.CreateNotification(notification); err != nil {
			return false, err
		}
	}

	return true, nil
}

func (s *Shoppinglist) Edit() error {
	shoppinglist := map[string]interface{}{
		"id":    s.ID,
		"title": s.Title,
		"items": s.Items,
		"owner": s.Owner,
	}

	return models.EditList(s.ID, shoppinglist)
}

func (s *Shoppinglist) GetList() (*models.Shoppinglist, error) {
	return models.GetList(s.ID, s.Owner)
}

func (s *Shoppinglist) GetListsByOwner() (*[]models.Shoppinglist, error) {
	return models.GetListByEmail(s.Owner)
}

func (s *Shoppinglist) SendInvitationEmails() error {
	for idx, val := range s.Participants {
		email := s.Participants[idx]
		log.Print(val)
		log.Print(email)
		//TODO: Send Emails
	}
	return nil
}

func (s *Shoppinglist) GetLastPosition() (int64, error) {
	return models.GetLastPosition(s.ID)
}

func (s *Shoppinglist) Delete() error {
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
