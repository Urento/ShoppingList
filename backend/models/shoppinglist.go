package models

import (
	"errors"

	"gorm.io/gorm/clause"
)

type Shoppinglist struct {
	Model

	ID           int            `json:"id" gorm:"primaryKey"`
	Title        string         `json:"title"`
	Items        []*Item        `json:"items" gorm:"foreignKey:ParentListID;"`
	Owner        string         `json:"owner"`
	Participants []*Participant `json:"participants" gorm:"foreignKey:ParentListID;"`
}

//TODO: Split Item add Participants in different files
type Item struct {
	Model

	ID           int    `gorm:"primaryKey" json:"id"`
	ParentListID int    `json:"parentListId"`
	ItemID       int    `json:"itemId"`
	Title        string `json:"title"`
	Position     int64  `json:"position"`
	Bought       bool   `json:"bought" gorm:"default:false"`
}

type Participant struct {
	Model

	ID           int    `gorm:"primaryKey" json:"id"`
	ParentListID int    `json:"parentListId"`
	Status       string `json:"status" gorm:"default:'pending'"`
	Email        string `json:"email"`
}

func ExistByID(id int) (bool, error) {
	var Found bool
	err := db.Raw("SELECT EXISTS(SELECT created_on FROM shoppinglists WHERE id = ?) AS found", id).Scan(&Found).Error
	return Found, err
}

func GetTotalListsByOwner(ownerID string) (int64, error) {
	var count int64
	if err := db.Model(&Shoppinglist{}).Where("owner = ?", ownerID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetLists(owner string) ([]Shoppinglist, error) {
	var lists []Shoppinglist
	err := db.Preload("Participants").Where("owner = ?", owner).Find(&lists).Error
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func GetList(id int, owner string) (*Shoppinglist, error) {
	var list Shoppinglist
	err := db.Preload(clause.Associations).Where("id = ?", id).Where("owner = ?", owner).First(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func GetListByEmail(email string) (*[]Shoppinglist, error) {
	var list []Shoppinglist
	err := db.Model(&Shoppinglist{}).Preload(clause.Associations).Where("owner = ?", email).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func EditList(id int, data map[string]interface{}) error {
	shoppinglist := Shoppinglist{
		ID:    data["id"].(int),
		Title: data["title"].(string),
		Owner: data["owner"].(string),
	}
	err := db.Omit(clause.Associations).Where("id = ?", id).Updates(&shoppinglist).Error
	return err
}

func CreateList(data map[string]interface{}) error {
	shoppinglist := Shoppinglist{
		ID:           data["id"].(int),
		Title:        data["title"].(string),
		Owner:        data["owner"].(string),
		Participants: data["participants"].([]*Participant),
	}

	if err := db.Omit("Items").Create(&shoppinglist).Error; err != nil {
		return err
	}
	return nil
}

func DeleteList(id int) error {
	itemsCount := db.Model(&Shoppinglist{}).Where("id = ?", id).Association("Items").Count()

	if itemsCount >= 1 {
		if err := db.Model(&Shoppinglist{}).Where("id = ?", id).Association("Items").Delete(&Shoppinglist{ID: id}); err != nil {
			return err
		}
	}

	if err := db.Where("id = ?", id).Delete(&Shoppinglist{ID: id}).Error; err != nil {
		return err
	}
	return nil
}

func AddItem(item Item) (*Item, error) {
	exists, err := ExistByID(item.ParentListID)
	if err != nil || !exists {
		return nil, errors.New("shoppinglist does not exist")
	}

	err = db.Create(&item).Error
	return &item, err
}

func UpdateItem(item Item) error {
	exists, err := ExistByID(item.ParentListID)
	if err != nil || !exists {
		return errors.New("shoppinglist does not exist")
	}

	err = db.Model(&Item{}).Where("parent_list_id = ?", item.ParentListID).Where("item_id = ?", item.ItemID).Updates(&item).Error

	return err
}

func GetItem(id, itemID int) (Item, error) {
	var item Item
	err := db.Model(&Item{}).Where("parent_list_id = ?", id).Where("item_id = ?", itemID).First(&item).Error
	return item, err
}

func GetItems(id int) ([]Item, error) {
	var Items []Item
	err := db.Where("parent_list_id = ?", id).Preload("Items").Find(&Items).Error
	return Items, err
}

func GetLastPosition(id int) (int64, error) {
	var Position int64
	//err := db.Where("parent_list_id = ?", id).Preload("Items").Select("position").First(&Position).Error
	err := db.Model(&Item{}).Select("position").Where("parent_list_id = ?", id).Order("position asc").Find(&Position).Error
	if err != nil {
		return 0, err
	}
	return Position, nil
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
