package models

import "errors"

type Participant struct {
	Model

	ID           int    `gorm:"primaryKey" json:"id"`
	ParentListID int    `json:"parentListId"`
	Status       string `json:"status" gorm:"default:'pending'"`
	Email        string `json:"email"`
}

func AddParticipant(participant Participant) (Participant, error) {
	exists, err := ExistByID(participant.ParentListID)
	if err != nil || !exists {
		return Participant{}, errors.New("shoppinglist does not exist")
	}

	err = db.Model(&Participant{}).Create(&participant).Error
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
