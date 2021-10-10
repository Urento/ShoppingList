package models

import (
	"errors"
)

type Notification struct {
	Model

	ID               int    `gorm:"primaryKey" json:"id"`
	UserID           int    `json:"userId"`
	NotificationType string `json:"notification_type"`
	Title            string `json:"title"`
	Text             string `json:"text"`
	Read             bool   `json:"read" gorm:"default:false"`
}

func CreateNotification(notification Notification) error {
	exists, err := ExistsUserID(notification.UserID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("user not found")
	}

	err = db.Create(&notification).Error
	return err
}

func HasUnreadNotifications(userId int) (bool, error) {
	exists, err := ExistsUserID(userId)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, errors.New("user not found")
	}

	var Count int64
	err = db.Model(&Notification{}).
		Where("user_id = ?", userId).Where("read = ?", false).
		Limit(1).Count(&Count).Error

	if err != nil {
		return false, err
	}

	if Count > 0 {
		return true, nil
	}

	return false, nil
}

func GetNotifications(userId int) ([]Notification, error) {
	exists, err := ExistsUserID(userId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("user not found")
	}

	var Notifications []Notification
	err = db.Model(&Notification{}).Where("user_id = ?", userId).Find(&Notifications).Error

	return Notifications, err
}

func GetNotification(userId, id int) (Notification, error) {
	exists, err := ExistsUserID(userId)
	if err != nil {
		return Notification{}, err
	}

	if !exists {
		return Notification{}, errors.New("user not found")
	}

	var notification Notification
	err = db.Model(&Notification{}).Where("user_id = ?", userId).Where("id = ?", id).Limit(1).Find(&notification).Error

	return notification, err
}

func DeleteNotification(userId, id int) error {
	exists, err := ExistsUserID(userId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("user not found")
	}

	err = db.Where("user_id = ?", userId).Where("id = ?", id).Delete(&Notification{ID: id, UserID: userId}).Error

	return err
}

func MarkNotificationAsRead(userId, id int) error {
	exists, err := ExistsUserID(userId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("user not found")
	}

	err = db.Model(&Notification{}).Where("user_id = ?", userId).Where("id = ?", id).Update("read", true).Error
	return err
}

func MarkAllNotificationsAsRead(userId int) error {
	exists, err := ExistsUserID(userId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("user not found")
	}

	var Notifications []Notification
	err = db.Model(&Notification{}).Where("user_id = ?", userId).Find(&Notifications).Error
	if err != nil {
		return err
	}

	tx := db.Begin()

	for _, notification := range Notifications {
		if err := tx.Model(&Notification{}).Where("user_id = ?", notification.UserID).Where("id = ?", notification.ID).Update("read", true).Error; err != nil {
			return err
		}
	}

	tx.Commit()

	return nil
}
