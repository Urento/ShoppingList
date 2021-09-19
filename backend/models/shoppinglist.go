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

type Item struct {
	Model

	ID           int    `gorm:"primaryKey" json:"id"`
	ParentListID int    `json:"parentListId"`
	ItemID       int    `json:"itemId"`
	Title        string `json:"title"`
	Position     int    `json:"position"`
	Bought       bool   `json:"bought" gorm:"default:false"`
}

type Participant struct {
	Model

	ID           int     `json:"participantId"`
	ParentListID int     `json:"parentListId"`
	Status       *string `json:"status" gorm:"default:'pending'"`
	Email        string  `json:"email"`
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
	err := db.Debug().Preload(clause.Associations).Where("id = ?", id).Where("owner = ?", owner).First(&list).Error
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
	err := db.Debug().Omit(clause.Associations).Where("id = ?", id).Updates(&shoppinglist).Error
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
	if err := db.Debug().Where("id = ?", id).Delete(&Shoppinglist{ID: id}).Error; err != nil {
		return err
	}
	return nil
}

func AddItem(item Item) (*Item, error) {
	exists, err := ExistByID(item.ParentListID)
	if err != nil || !exists {
		return nil, errors.New("shoppinglist does not exist")
	}

	//err = db.Debug().Create(&item).Error
	err = db.Debug().Model(&Shoppinglist{}).Association("Items").Append(item)

	return &item, err
}

func UpdateItem(id, itemID int) error {
	return nil
}

func GetItem(id, itemID int) (Item, error) {
	i := Item{}
	return i, nil
}

func GetItems(id int) ([]Item, error) {
	var items []Item
	err := db.Debug().Where("parent_list_id = ?", id).Preload("Items").Find(&items).Error
	return items, err
}

func GetLastPosition(id int) (int64, error) {
	var Position int64
	err := db.Debug().Where("parent_list_id = ?", id).Preload("Items").Select("position").Order("position desc").First(&Position).Error
	return Position, err
}
