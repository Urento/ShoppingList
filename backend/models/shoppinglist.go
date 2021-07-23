package models

import (
	"fmt"

	"github.com/lib/pq"
)

type Shoppinglist struct {
	Model
	ID           int `gorm:"primaryKey"`
	Title        string
	Items        pq.StringArray `gorm:"type:text[]"`
	Owner        string
	Participants pq.StringArray `gorm:"type:text[]"`
	Position     int
}

func ExistByID(id int) (bool, error) {
	var Found bool
	err := db.Raw("SELECT EXISTS(SELECT created_on FROM shoppinglists WHERE id = ? AND deleted_at = ?) AS found", id, nil).Scan(&Found).Error
	return Found, err
}

func GetTotalListsByOwner(ownerID string) (int64, error) {
	var count int64
	if err := db.Model(&Shoppinglist{}).Where("owner = ? AND deleted_at = ?", ownerID, nil).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetLists(owner string) ([]Shoppinglist, error) {
	var lists []Shoppinglist
	err := db.Raw("SELECT * FROM shoppinglists WHERE owner = ? AND deleted_at = ?", owner, nil).Scan(&lists).Error
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func GetList(id int) (*Shoppinglist, error) {
	var list Shoppinglist
	err := db.Model(&Shoppinglist{}).Where("id = ?", id).First(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func GetListByEmail(email string) (*[]Shoppinglist, error) {
	var list []Shoppinglist
	err := db.Model(&Shoppinglist{}).Where("owner = ? AND deleted_at = ?", email, nil).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func EditList(id int, data map[string]interface{}) error {
	shoppinglist := Shoppinglist{
		Title:        data["title"].(string),
		Items:        data["items"].([]string),
		Owner:        data["owner"].(string),
		Participants: data["participants"].([]string),
		Position:     data["position"].(int),
	}

	if err := db.Model(&Shoppinglist{}).Where("id = ?", id).Updates(shoppinglist).Error; err != nil {
		return err
	}

	return nil
}

func CreateList(data map[string]interface{}) error {
	shoppinglist := Shoppinglist{
		ID:           data["id"].(int),
		Title:        data["title"].(string),
		Items:        data["items"].([]string),
		Owner:        data["owner"].(string),
		Participants: data["participants"].([]string),
		Position:     data["position"].(int),
	}

	if err := db.Create(&shoppinglist).Error; err != nil {
		return err
	}
	return nil
}

func DeleteList(id int) error {
	if err := db.Where("id = ?", id).Delete(&Shoppinglist{}).Error; err != nil {
		return err
	}
	return nil
}

func DropTable() {
	db.Migrator().DropTable(&Shoppinglist{})
	fmt.Println("Dropping Shoppinglist Table")
}
