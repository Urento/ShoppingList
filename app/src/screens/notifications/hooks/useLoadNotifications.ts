import { useState } from "react-router/node_modules/@types/react";

const useLoadNotifications = () => {
  const [loadingNotification, setLoadingNotifications] =
    useState<boolean>(false);
  const [notifications, setNotifications] = useState<any>([]);
};

export default useLoadNotifications;
