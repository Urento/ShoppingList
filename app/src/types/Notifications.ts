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
  id?: string;
  userId?: string;
  notification_type: string;
  title: string;
  text: string;
  read?: boolean;
}
