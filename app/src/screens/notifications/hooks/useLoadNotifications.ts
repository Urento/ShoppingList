import { useEffect, useState } from "react";
import { NotificationsResponse } from "../../../types/Notifications";
import { API_URL } from "../../../util/constants";

const useLoadNotifications = () => {
  const [loadingNotification, setLoadingNotifications] =
    useState<boolean>(true);
  const [notifications, setNotifications] = useState<any>([]);

  const loadNotifications = async () => {
    const response = await fetch(`${API_URL}/notifications`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    });
    const fJson: NotificationsResponse = await response.json();
    setNotifications(fJson.data);
    setLoadingNotifications(false);
  };

  useEffect(() => {
    loadNotifications();
  }, []);

  return [loadingNotification, notifications];
};

export default useLoadNotifications;
