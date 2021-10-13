import { useEffect, useState } from "react";
import { HasUnreadNotificationsResponse } from "../types/Notifications";
import { API_URL } from "../util/constants";

const useHasUnreadNotifications = () => {
  const [hasUnreadNotifications, setHasUnreadNotifications] =
    useState<boolean>(false);

  useEffect(() => {
    const checkNotifications = async () => {
      const response = await fetch(`${API_URL}/notifications/n/hasunread`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        credentials: "include",
      });
      const fJson: HasUnreadNotificationsResponse = await response.json();
      if (fJson.data.success === "true")
        setHasUnreadNotifications(fJson.data.has === "true");
    };

    checkNotifications();
  }, []);

  return hasUnreadNotifications;
};

export default useHasUnreadNotifications;
