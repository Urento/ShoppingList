package models

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/pkg/util"
)

func CreateUser() (int, error) {
	username := util.StringWithCharset(500)
	email := util.RandomEmail()
	password := util.StringWithCharset(500)
	ip := util.RandomIPAddress()

	err := CreateAccount(email, username, password, ip)
	if err != nil {
		return 0, err
	}

	user, err := GetUser(email)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func TestCreateNotification(t *testing.T) {
	Setup()

	userId, err := CreateUser()
	if err != nil {
		t.Errorf("Error while creating user: %s", err)
	}

	notificationType := "invitation"
	text := util.StringWithCharset(500)
	title := util.StringWithCharset(300)

	notification := Notification{
		UserID:           userId,
		Title:            title,
		NotificationType: notificationType,
		Text:             text,
		Read:             false,
	}

	err = CreateNotification(notification)
	if err != nil {
		t.Errorf("Error while creating notification: %s", err)
	}

	Equal(t, nil, err)
}

func TestHasUnreadNotifications(t *testing.T) {
	Setup()

	t.Run("Has unread notifications", func(t *testing.T) {
		userId, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		notificationType := "invitation"
		text := util.StringWithCharset(500)
		title := util.StringWithCharset(300)

		notification := Notification{
			UserID:           userId,
			Title:            title,
			NotificationType: notificationType,
			Text:             text,
			Read:             false,
		}

		err = CreateNotification(notification)
		if err != nil {
			t.Errorf("Error while creating notification: %s", err)
		}

		has, err := HasUnreadNotifications(userId)
		if err != nil {
			t.Errorf("Error while getting unread notifications: %s", err)
		}

		Equal(t, true, has)
	})

	t.Run("Has unread notification but he has no unread notifications", func(t *testing.T) {
		userId, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		has, err := HasUnreadNotifications(userId)
		if err != nil {
			t.Errorf("Error while getting unread notifications: %s", err)
		}

		Equal(t, false, has)
	})
}

func TestGetNotifications(t *testing.T) {
	Setup()

	t.Run("Get Notifications", func(t *testing.T) {
		userId, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		notificationType := "invitation"
		text := util.StringWithCharset(500)
		title := util.StringWithCharset(300)

		notification := Notification{
			UserID:           userId,
			Title:            title,
			NotificationType: notificationType,
			Text:             text,
			Read:             false,
		}

		err = CreateNotification(notification)
		if err != nil {
			t.Errorf("Error while creating notification: %s", err)
		}

		notifications, err := GetNotifications(userId)
		if err != nil {
			t.Errorf("Error while getting notifications: %s", err)
		}

		Equal(t, title, notifications[0].Title)
		Equal(t, text, notifications[0].Text)
		Equal(t, notificationType, notifications[0].NotificationType)
		Equal(t, false, notifications[0].Read)
		Equal(t, userId, notifications[0].UserID)
	})

	t.Run("Get Notifications when the user has no notifications", func(t *testing.T) {
		userId, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		notifications, err := GetNotifications(userId)
		if err != nil {
			t.Errorf("Error while getting notifications: %s", err)
		}

		if len(notifications) > 0 {
			t.Errorf("Notifications are not an empty array")
		}
	})
}

func TestGetNotification(t *testing.T) {
	Setup()

	t.Run("Get Notification", func(t *testing.T) {
		userId, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		notificationType := "invitation"
		text := util.StringWithCharset(500)
		title := util.StringWithCharset(300)

		notification := Notification{
			UserID:           userId,
			Title:            title,
			NotificationType: notificationType,
			Text:             text,
			Read:             false,
		}

		err = CreateNotification(notification)
		if err != nil {
			t.Errorf("Error while creating notification: %s", err)
		}

		notifications, err := GetNotifications(userId)
		if err != nil {
			t.Errorf("Error while getting notifications: %s", err)
		}

		n, err := GetNotification(userId, notifications[0].ID)
		if err != nil {
			t.Errorf("Error while getting notification: %s", err)
		}

		Equal(t, title, n.Title)
		Equal(t, text, n.Text)
		Equal(t, userId, n.UserID)
		Equal(t, false, n.Read)
		Equal(t, notificationType, n.NotificationType)
	})

	t.Run("Get Notification that doesn't exist", func(t *testing.T) {
		userId, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		_, err = GetNotification(userId, 456784896324532)

		Equal(t, nil, err)
	})
}

func TestDeleteNotification(t *testing.T) {
	Setup()

	t.Run("Delete Notification", func(t *testing.T) {
		userId, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		notificationType := "invitation"
		text := util.StringWithCharset(500)
		title := util.StringWithCharset(300)

		notification := Notification{
			UserID:           userId,
			Title:            title,
			NotificationType: notificationType,
			Text:             text,
			Read:             false,
		}

		err = CreateNotification(notification)
		if err != nil {
			t.Errorf("Error while creating notification: %s", err)
		}

		notifications, err := GetNotifications(userId)
		if err != nil {
			t.Errorf("Error while getting notifications: %s", err)
		}

		err = DeleteNotification(userId, notifications[0].ID)
		if err != nil {
			t.Errorf("Error while deleting notification: %s", err)
		}

		notificationsAfter, err := GetNotifications(userId)
		if err != nil {
			t.Errorf("Error while getting notifications: %s", err)
		}

		if len(notificationsAfter) > 0 {
			t.Errorf("Notification did not get deleted")
		}
	})

	t.Run("Delete Notification that doesn't exist", func(t *testing.T) {
		userId, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		err = DeleteNotification(userId, 456784896324532)

		Equal(t, nil, err)
	})
}
