package models

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/pkg/util"
)

//TODO: Add Test where I check if the notification was created after the shoppinglist was edited, deleted and created

func CreateUser() (*Auth, error) {
	username := util.StringWithCharset(500)
	email := util.RandomEmail()
	password := util.StringWithCharset(500)
	ip := util.RandomIPAddress()

	err := CreateAccount(email, username, password, ip)
	if err != nil {
		return nil, err
	}

	user, err := GetUser(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func TestCreateNotification(t *testing.T) {
	Setup()

	user, err := CreateUser()
	if err != nil {
		t.Errorf("Error while creating user: %s", err)
	}

	notificationType := "invitation"
	text := util.StringWithCharset(500)
	title := util.StringWithCharset(300)

	notification := Notification{
		UserID:           user.ID,
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
		user, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		notificationType := "invitation"
		text := util.StringWithCharset(500)
		title := util.StringWithCharset(300)

		notification := Notification{
			UserID:           user.ID,
			Title:            title,
			NotificationType: notificationType,
			Text:             text,
			Read:             false,
		}

		err = CreateNotification(notification)
		if err != nil {
			t.Errorf("Error while creating notification: %s", err)
		}

		has, err := HasUnreadNotifications(user.ID)
		if err != nil {
			t.Errorf("Error while getting unread notifications: %s", err)
		}

		Equal(t, true, has)
	})

	t.Run("Has unread notification but he has no unread notifications", func(t *testing.T) {
		user, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		has, err := HasUnreadNotifications(user.ID)
		if err != nil {
			t.Errorf("Error while getting unread notifications: %s", err)
		}

		Equal(t, false, has)
	})
}

func TestGetNotifications(t *testing.T) {
	Setup()

	t.Run("Get Notifications", func(t *testing.T) {
		user, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		notificationType := "invitation"
		text := util.StringWithCharset(500)
		title := util.StringWithCharset(300)

		notification := Notification{
			UserID:           user.ID,
			Title:            title,
			NotificationType: notificationType,
			Text:             text,
			Read:             false,
		}

		err = CreateNotification(notification)
		if err != nil {
			t.Errorf("Error while creating notification: %s", err)
		}

		notifications, err := GetNotifications(user.ID)
		if err != nil {
			t.Errorf("Error while getting notifications: %s", err)
		}

		Equal(t, title, notifications[0].Title)
		Equal(t, text, notifications[0].Text)
		Equal(t, notificationType, notifications[0].NotificationType)
		Equal(t, false, notifications[0].Read)
		Equal(t, user.ID, notifications[0].UserID)
	})

	t.Run("Get Notifications when the user has no notifications", func(t *testing.T) {
		user, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		notifications, err := GetNotifications(user.ID)
		if err != nil {
			t.Errorf("Error while getting notifications: %s", err)
		}

		if len(notifications) > 0 {
			t.Errorf("Notifications are not an empty array")
		}
	})

	t.Run("Get Notifications with more then 50 notifications", func(t *testing.T) {
		user, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		for i := 0; i < 60; i++ {
			notificationType := "invitation"
			text := util.StringWithCharset(5000)
			title := util.StringWithCharset(3000)

			notification := Notification{
				UserID:           user.ID,
				Title:            title,
				NotificationType: notificationType,
				Text:             text,
				Read:             false,
			}

			err = CreateNotification(notification)
			if err != nil {
				t.Errorf("Error while creating notification: %s", err)
			}
		}

		notifications, err := GetNotifications(user.ID)
		if err != nil {
			t.Errorf("Error while getting notifications: %s", err)
		}

		if len(notifications) > 50 {
			t.Errorf("More then 50 notifications were loaded")
		}
	})
}

func TestGetNotification(t *testing.T) {
	Setup()

	t.Run("Get Notification", func(t *testing.T) {
		user, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		notificationType := "invitation"
		text := util.StringWithCharset(500)
		title := util.StringWithCharset(300)

		notification := Notification{
			UserID:           user.ID,
			Title:            title,
			NotificationType: notificationType,
			Text:             text,
			Read:             false,
		}

		err = CreateNotification(notification)
		if err != nil {
			t.Errorf("Error while creating notification: %s", err)
		}

		notifications, err := GetNotifications(user.ID)
		if err != nil {
			t.Errorf("Error while getting notifications: %s", err)
		}

		n, err := GetNotification(user.ID, notifications[0].ID)
		if err != nil {
			t.Errorf("Error while getting notification: %s", err)
		}

		Equal(t, title, n.Title)
		Equal(t, text, n.Text)
		Equal(t, user.ID, n.UserID)
		Equal(t, false, n.Read)
		Equal(t, notificationType, n.NotificationType)
	})

	t.Run("Get Notification that doesn't exist", func(t *testing.T) {
		user, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		_, err = GetNotification(user.ID, 456784896324532)

		Equal(t, nil, err)
	})
}

func TestDeleteNotification(t *testing.T) {
	Setup()

	t.Run("Delete Notification", func(t *testing.T) {
		user, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		notificationType := "invitation"
		text := util.StringWithCharset(500)
		title := util.StringWithCharset(300)

		notification := Notification{
			UserID:           user.ID,
			Title:            title,
			NotificationType: notificationType,
			Text:             text,
			Read:             false,
		}

		err = CreateNotification(notification)
		if err != nil {
			t.Errorf("Error while creating notification: %s", err)
		}

		notifications, err := GetNotifications(user.ID)
		if err != nil {
			t.Errorf("Error while getting notifications: %s", err)
		}

		err = DeleteNotification(user.ID, notifications[0].ID)
		if err != nil {
			t.Errorf("Error while deleting notification: %s", err)
		}

		notificationsAfter, err := GetNotifications(user.ID)
		if err != nil {
			t.Errorf("Error while getting notifications: %s", err)
		}

		if len(notificationsAfter) > 0 {
			t.Errorf("Notification did not get deleted")
		}
	})

	t.Run("Delete Notification that doesn't exist", func(t *testing.T) {
		user, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		err = DeleteNotification(user.ID, 456784896324532)

		Equal(t, nil, err)
	})
}

func TestMarkNotificationAsRead(t *testing.T) {
	Setup()

	user, err := CreateUser()
	if err != nil {
		t.Errorf("Error while creating user: %s", err)
	}

	notificationType := "invitation"
	text := util.StringWithCharset(500)
	title := util.StringWithCharset(300)

	notification := Notification{
		UserID:           user.ID,
		Title:            title,
		NotificationType: notificationType,
		Text:             text,
		Read:             false,
	}

	err = CreateNotification(notification)
	if err != nil {
		t.Errorf("Error while creating notification: %s", err)
	}

	notifications, err := GetNotifications(user.ID)
	if err != nil {
		t.Errorf("Error while getting notifications: %s", err)
	}

	err = MarkNotificationAsRead(user.ID, notifications[0].ID)
	if err != nil {
		t.Errorf("Error while marking a notification as read: %s", err)
	}

	notificationsAfter, err := GetNotifications(user.ID)
	if err != nil {
		t.Errorf("Error while getting notifications: %s", err)
	}

	Equal(t, true, notificationsAfter[0].Read)
}

func TestMarkAllNotificationsAsRead(t *testing.T) {
	Setup()

	user, err := CreateUser()
	if err != nil {
		t.Errorf("Error while creating user: %s", err)
	}

	notificationType := "invitation"
	text := util.StringWithCharset(500)
	title := util.StringWithCharset(300)

	notification := Notification{
		UserID:           user.ID,
		Title:            title,
		NotificationType: notificationType,
		Text:             text,
		Read:             false,
	}

	err = CreateNotification(notification)
	if err != nil {
		t.Errorf("Error while creating notification: %s", err)
	}

	text2 := util.StringWithCharset(500)
	title2 := util.StringWithCharset(300)

	notification2 := Notification{
		UserID:           user.ID,
		Title:            title2,
		NotificationType: notificationType,
		Text:             text2,
		Read:             false,
	}

	err = CreateNotification(notification2)
	if err != nil {
		t.Errorf("Error while creating notification: %s", err)
	}

	err = MarkAllNotificationsAsRead(user.ID)
	if err != nil {
		t.Errorf("Error while marking all notifications as read: %s", err)
	}

	notifications, err := GetNotifications(user.ID)
	if err != nil {
		t.Errorf("Error while getting notifications: %s", err)
	}

	Equal(t, true, notifications[0].Read)
	Equal(t, true, notifications[1].Read)
}
