package models

import (
	"errors"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Shoppinglist struct {
	Model

	ID           int            `gorm:"primaryKey" json:"id"`
	Title        string         `json:"title"`
	Items        []*Item        `json:"items" gorm:"foreignKey:ParentListID;"`
	Owner        string         `json:"owner"`
	Participants pq.StringArray `gorm:"type:text[]" json:"participants"`
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

//TODO: Implement Participant to Shoppinglist Struct
type Participant struct {
	Model

	ID     int    `json:"participantId"`
	Status string `json:"status"`
	Email  string `json:"email"`
}

func ExistByID(id int) (bool, error) {
	var Found bool
	err := db.Raw("SELECT EXISTS(SELECT created_on FROM shoppinglists WHERE id = ?) AS found", id).Scan(&Found).Error
	return Found, err
}

func GetTotalListsByOwner(ownerID string) (int64, error) {
	var count int64
	if err := db.Model(&Shoppinglist{}).Preload("Items").Where("owner = ?", ownerID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetLists(owner string) ([]Shoppinglist, error) {
	var lists []Shoppinglist
	err := db.Preload("Items").Where("owner = ?", owner).Find(&lists).Error
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func GetList(id int, owner string) (*Shoppinglist, error) {
	var list Shoppinglist
	err := db.Debug().Preload("Items").Where("id = ? AND owner = ?", id, owner).First(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func GetListByEmail(email string) (*[]Shoppinglist, error) {
	var list []Shoppinglist
	err := db.Model(&Shoppinglist{}).Preload("Items").Where("owner = ?", email).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func EditList(id int, data map[string]interface{}) error {
	shoppinglist := Shoppinglist{
		ID:           data["id"].(int),
		Title:        data["title"].(string),
		Items:        data["items"].([]*Item),
		Owner:        data["owner"].(string),
		Participants: data["participants"].([]string),
	}
	err := db.Debug().Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", id).Updates(&shoppinglist).Error
	return err
}

func AddItem(item Item) error {
	exists, err := ExistByID(item.ParentListID)
	if err != nil || !exists {
		return errors.New("shoppinglist does not exist")
	}

	err = db.Debug().Create(&item).Error

	return err
}

func UpdateItem(itemId int) error {
	return nil
}

func GetItem(itemId int) (Item, error) {
	i := Item{}
	return i, nil
}

func CreateList(data map[string]interface{}) error {
	shoppinglist := Shoppinglist{
		ID:           data["id"].(int),
		Title:        data["title"].(string),
		Owner:        data["owner"].(string),
		Participants: data["participants"].([]string),
	}

	if err := db.Debug().Omit(clause.Associations).Create(&shoppinglist).Error; err != nil {
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

func GetItems(id int, owner string) ([]Item, error) {
	var items []Item
	//err := db.Model(&Shoppinglist{}).Where("id = ? AND owner = ?", id, owner).Association("Items").Find(&items).Error()
	err := db.Model(&Shoppinglist{}).Where("id = ?", id, owner).Association("Items").Find(&items).Error()
	return items, errors.New(err)
}
