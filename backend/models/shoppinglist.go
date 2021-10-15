package models

import (
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

func CreateList(data Shoppinglist) error {
	shoppinglist := Shoppinglist{
		ID:    data.ID,
		Title: data.Title,
		Owner: data.Owner,
	}

	err := db.Debug().Model(&Shoppinglist{}).Omit(clause.Associations).Create(&shoppinglist).Error
	return err
}

func DeleteList(id int) error {
	itemsCount := db.Model(&Shoppinglist{}).Where("id = ?", id).Association("Items").Count()

	if itemsCount > 0 {
		if err := db.Model(&Shoppinglist{}).Where("id = ?", id).Association("Items").Delete(&Shoppinglist{ID: id}); err != nil {
			return err
		}
	}

	err := db.Where("id = ?", id).Delete(&Shoppinglist{ID: id}).Error
	return err
}
