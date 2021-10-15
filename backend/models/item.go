package models

import "errors"

type Item struct {
	Model

	ID           int    `gorm:"primaryKey" json:"id"`
	ParentListID int    `json:"parentListId"`
	ItemID       int    `json:"itemId"`
	Title        string `json:"title"`
	Position     int64  `json:"position"`
	Bought       bool   `json:"bought" gorm:"default:false"`
}

func AddItem(item Item) (*Item, error) {
	exists, err := ExistByID(item.ParentListID)
	if err != nil || !exists {
		return nil, errors.New("shoppinglist does not exist")
	}

	err = db.Create(&item).Error
	return &item, err
}

func DeleteItem(parentListId, id int) error {
	exists, err := ExistByID(parentListId)
	if err != nil || !exists {
		return errors.New("shoppinglist does not exist")
	}

	err = db.Model(&Item{}).Where("item_id = ?", id).Where("parent_list_id = ?", parentListId).Delete(&Item{}).Error
	return err
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
	//err := db.Where("parent_list_id = ?", id).Preload("Items").Find(&Items).Error
	err := db.Model(&Item{}).Where("parent_list_id = ?", id).Find(&Items).Error
	return Items, err
}

func GetLastPosition(id int) (int64, error) {
	var Position int64
	err := db.Model(&Item{}).Select("position").Where("parent_list_id = ?", id).Order("position asc").Find(&Position).Error
	if err != nil {
		return 0, err
	}
	return Position, nil
}
