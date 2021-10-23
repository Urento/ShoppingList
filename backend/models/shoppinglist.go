package models

import (
	"fmt"
	"time"

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

func GetLists(owner string, offset int) ([]Shoppinglist, error) {
	var lists []Shoppinglist
	err := db.Preload("Participants").Omit("Items").Where("owner = ?", owner).Limit(6).Offset(offset).Find(&lists).Error
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

func GetListWithoutOwner(id int) (*Shoppinglist, error) {
	var list Shoppinglist
	err := db.Preload(clause.Associations).Where("id = ?", id).First(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func GetListByEmail(email string, offset int) (*[]Shoppinglist, error) {
	var list []Shoppinglist
	err := db.Model(&Shoppinglist{}).Preload("Participants").Where("owner = ?", email).Limit(6).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func EditList(id int, data Shoppinglist) error {
	err := db.Omit(clause.Associations).Where("id = ?", id).Updates(&data).Error
	return err
}

func CreateList(data Shoppinglist, userId int, withNotification bool) error {
	if withNotification {
		notification := Notification{
			UserID:           userId,
			Title:            "New Shoppinglist",
			Text:             fmt.Sprintf("%s was created", data.Title),
			NotificationType: "new_shoppinglist",
			Date:             time.Now().Format("02.01.2006"),
		}

		if err := CreateNotification(notification); err != nil {
			return err
		}
	}

	err := db.Model(&Shoppinglist{}).Omit(clause.Associations).Create(&data).Error
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

func BelongsShoppinglistToEmail(email string, id int) (bool, error) {
	var Count int64
	err := db.Model(&Shoppinglist{}).Where("id = ?", id).Where("owner = ?", email).Count(&Count).Limit(1).Error
	return Count >= 1, err
}
