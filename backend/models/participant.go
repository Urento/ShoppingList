package models

import (
	"errors"
)

type Participant struct {
	Model

	ID           int    `gorm:"primaryKey" json:"id"`
	ParentListID int    `json:"parentListId"`
	Status       string `json:"status" gorm:"default:'pending'"`
	Email        string `json:"email"`
	RequestFrom  string `json:"request_from"`
}

func AddParticipant(participant Participant) (Participant, error) {
	exists, err := ExistByID(participant.ParentListID)
	if err != nil || !exists {
		return Participant{}, errors.New("shoppinglist does not exist")
	}

	err = db.Create(&participant).Error
	return participant, err
}

func RemoveParticipant(parentListID, id int) error {
	exists, err := ExistByID(parentListID)
	if err != nil || !exists {
		return errors.New("shoppinglist does not exist")
	}

	err = db.Model(&Participant{}).Where("parent_list_id = ?", parentListID).Where("id = ?", id).Delete(&Participant{}).Error
	return err
}

func GetParticipants(parentListID int) ([]Participant, error) {
	exists, err := ExistByID(parentListID)
	if err != nil || !exists {
		return nil, errors.New("shoppinglist does not exist")
	}

	var Participants []Participant
	err = db.Model(&Participant{}).Where("parent_list_id = ?", parentListID).Find(&Participants).Error
	return Participants, err
}

func IsParticipantAlreadyIncluded(email string, parentListID int) (bool, error) {
	var Count int64
	err := db.Model(&Participant{}).Where("parent_list_id = ?", parentListID).Where("email = ?", email).Limit(1).Count(&Count).Error
	return Count >= 1, err
}

func GetListsByParticipant(participantEmail string) ([]Shoppinglist, error) {
	listsByParticipants := []Participant{}
	lists := []Shoppinglist{}
	err := db.Model(&Participant{}).Where("email = ?", participantEmail).Where("status = ?", "accepted").Find(&listsByParticipants).Error
	if err != nil {
		return []Shoppinglist{}, nil
	}

	if len(listsByParticipants) <= 0 {
		return []Shoppinglist{}, nil
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, val := range listsByParticipants {
		var l Shoppinglist
		err = tx.Model(&Shoppinglist{}).Preload("Participants").Where("id = ?", val.ParentListID).Limit(1).Find(&l).Error
		if err != nil {
			return lists, err
		}
		lists = append(lists, l)
	}

	err = tx.Commit().Error
	if err != nil {
		return []Shoppinglist{}, err
	}

	return lists, nil
}

func GetPendingRequests(email string) ([]Participant, error) {
	var requests []Participant
	err := db.Model(&Participant{}).Where("status = ?", "pending").Where("email = ?", email).Find(&requests).Error
	if err != nil {
		return []Participant{}, err
	}
	return requests, nil
}

func GetPendingRequestsFromShoppinglist(email string, id int) ([]Participant, error) {
	belongs, err := BelongsShoppinglistToEmail(email, id)
	if err != nil {
		return []Participant{}, err
	}

	if !belongs {
		return []Participant{}, errors.New("shoppinglist doesn not belongs to the given email")
	}

	var requests []Participant
	err = db.Model(&Participant{}).Where("status = ?", "pending").Where("parent_list_id = ?", id).Find(&requests).Error
	if err != nil {
		return []Participant{}, err
	}
	return requests, nil
}

func AcceptRequest(id int, email string) error {
	tx := db.Begin()

	err := tx.Model(&Participant{}).Where("id = ?", id).Where("email = ?", email).Update("status", "accepted").Error
	if err != nil {
		return err
	}
	err = tx.Model(&Participant{}).Where("id = ?", id).Where("email = ?", email).Update("request_from", nil).Error
	if err != nil {
		return err
	}

	tx.Commit()
	return err
}

func DeleteRequest(id int, email string) error {
	err := db.Model(&Participant{}).Where("id = ?", id).Where("email = ?", email).Delete(&Participant{ID: id, Email: email}).Error
	return err
}

func DeleteAll(email string) error {
	var requests []Participant
	err := db.Model(&Participant{}).Where("email = ?", email).Find(&requests).Error
	if err != nil {
		return err
	}

	tx := db.Begin()

	for _, request := range requests {
		err = tx.Model(&Participant{}).Where("id = ?", request.ID).Where("email = ?", request.Email).Delete(&Participant{Email: email}).Error
		if err != nil {
			return err
		}
	}

	err = tx.Commit().Error
	return err
}

func LeaveShoppinglist(id int, email string) error {
	err := db.Model(&Participant{}).Where("parent_list_id = ?", id).Where("email = ?", email).Where("status = ?", "accepted").Delete(&Participant{ParentListID: id, Email: email}).Error
	return err
}
