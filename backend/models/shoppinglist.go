package models

import (
	"fmt"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Shoppinglist struct {
	Model
	ID           int `gorm:"primaryKey"`
	Title        string
	Items        pq.StringArray `gorm:"type:text[]"`
	Owner        string
	Participants pq.StringArray `gorm:"type:text[]"`
}

func ExistByID(id int) (bool, error) {
	count := int64(0)
	err := db.Model(&Shoppinglist{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func GetTotalListsByOwner(ownerID string) (int64, error) {
	var count int64
	if err := db.Model(&Shoppinglist{}).Where("owner = ?", ownerID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetLists(pageNum int, pageSize int, owner string) ([]*Shoppinglist, error) {
	var shoppinglists []*Shoppinglist
	err := db.Model(&Shoppinglist{}).Where("owner = ?", owner).Offset(pageNum).Limit(pageSize).Find(&shoppinglists).Error
	if err != nil {
		return nil, err
	}
	return shoppinglists, nil
}

func GetList(id int) (*Shoppinglist, error) {
	var list Shoppinglist
	err := db.Where("id = ?", id).First(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &list, nil
}

func EditList(id int, data interface{}) error {
	if err := db.Model(&Shoppinglist{}).Where("id = ?", id).Updates(data).Error; err != nil {
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
