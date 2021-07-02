package models

import "gorm.io/gorm"

type Shoppinglist struct {
	Model

	ListID int `json:"list_id" gorm:"index"`

	Title        string
	Items        []string
	Owner        string
	Participants []string
}

func ExistByID(id int) (bool, error) {
	var shoppinglist Shoppinglist
	err := db.Select("id").Where("id = ?", id).First(&shoppinglist).Error
	if err != nil {
		return false, err
	}
	if shoppinglist.ID > 0 {
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

func GetLists(pageNum int, pageSize int, ownerID string) ([]*Shoppinglist, error) {
	var shoppinglists []*Shoppinglist
	err := db.Model(&Shoppinglist{}).Where("owner = ?", ownerID).Offset(pageNum).Limit(pageSize).Find(&shoppinglists).Error
	if err != nil {
		return nil, err
	}
	return shoppinglists, nil
}

func GetList(id int) (*Shoppinglist, error) {
	var list Shoppinglist
	err := db.Where("id ? and deleted_on = ?", id, 0).First(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &list, nil
}

func EditList(id int, data interface{}) error {
	if err := db.Model(&Shoppinglist{}).Where("id = ? AND deleted_on = ?", id, 0).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func CreateList(data map[string]interface{}) error {
	list := Shoppinglist{
		ListID:       data["list_id"].(int),
		Title:        data["title"].(string),
		Items:        data["items"].([]string),
		Owner:        data["owner"].(string),
		Participants: data["participants"].([]string),
	}
	if err := db.Create(&list).Error; err != nil {
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
