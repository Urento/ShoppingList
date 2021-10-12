export interface NotificationsResponse {
  code: string;
  message: string;
  data: Notification[];
}

export interface NotificationResponse {
  code: string;
  message: string;
  data: Notification;
}

export interface HasUnreadNotificationsResponse {
  code: string;
  message: string;
  data: HasUnreadNotificationDataResponse;
}

interface HasUnreadNotificationDataResponse {
  success: "true" | "false";
  has: "true" | "false";
}

export interface Notification {
  created_on?: number;
  date?: string;
  id?: number;
  userId?: string;
  notification_type: string;
  title: string;
  text: string;
  read?: boolean;
}
