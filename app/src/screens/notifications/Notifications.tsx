import { useHistory } from "react-router";
import { useState } from "react";
import { Button } from "../../components/Button";
import { Loading } from "../../components/Loading";
import { Sidebar } from "../../components/Sidebar";
import useAuthCheck from "../../hooks/useAuthCheck";
import useLoadNotifications from "./hooks/useLoadNotifications";
import { Notification, NotificationResponse } from "../../types/Notifications";
import { API_URL } from "../../util/constants";
import swal from "sweetalert";

const Notifications: React.FC = ({}) => {
  const history = useHistory();
  const authStatus = useAuthCheck();
  const [deleting, setDeleting] = useState<boolean>(false);
  const [loadingNotification, setLoadingNotification] =
    useState<boolean>(false);
  const [markingAllNotificationsAsRead, setMarkingAllNotificationsAsRead] =
    useState<boolean>(false);
  const [loadingNotifications, notifications] = useLoadNotifications();

  if (authStatus === "fail") {
    localStorage.removeItem("authenticated");
    history.push("/");
  }

  if (authStatus === "pending") return <Loading withSidebar />;
  if (loadingNotifications) return <Loading withSidebar />;

  const deleteNotification = async (id: number) => {
    setDeleting(true);

    const alert = await swal({
      icon: "warning",
      title: "Are you sure you want to delete the notification?",
      buttons: ["No, dont delete!", "Yes, delete!"],
      dangerMode: true,
    });

    if (alert) {
      await fetch(`${API_URL}/notifications`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: JSON.stringify({ notification_id: id }),
        credentials: "include",
      });
      window.location.reload();
    }

    setDeleting(false);
  };

  const loadNotification = async (id: number): Promise<Notification> => {
    const response = await fetch(`${API_URL}/notification/${id}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    });
    const fJson: NotificationResponse = await response.json();
    setLoadingNotification(false);
    return fJson.data;
  };

  const viewNotification = async (id: number) => {
    setLoadingNotification(true);
    swal({
      title: "Loading Notification...",
    });

    const notification = await loadNotification(id);

    if (!loadingNotification) {
      swal.close!();
      swal({
        title: notification.title,
        text: notification.text,
        buttons: [false, "Close"],
      });
    }
  };

  const markAllNotificationsAsRead = async () => {
    console.log("hdfsgjhbdfg");
    setMarkingAllNotificationsAsRead(true);
    await fetch(`${API_URL}/notifications/n/markall`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    });
    setMarkingAllNotificationsAsRead(false);
  };

  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
      <div className="container mx-auto py-10 md:w-4/5 w-11/12">
        <div className="bg-white px-4 md:px-10 pt-5 md:pt-7 pb-5 overflow-y-auto">
          {notifications.length > 0 && (
            <Button
              text="Mark all Notifications as read"
              loadingText="Marking all notifications as read..."
              onClick={markAllNotificationsAsRead}
              className="inline-flex sm:ml-3 mt-4 sm:mt-0 items-start justify-start px-6 py-3 bg-green-600 hover:bg-green-500 text-white focus:outline-none rounded"
              loading={markingAllNotificationsAsRead}
              color="green"
            />
          )}
          {notifications.length <= 0 && <p>You have no notifications!</p>}
          <table className="w-full whitespace-nowrap">
            <thead>
              <tr className="h-16 w-full text-sm leading-none ">
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
              </tr>
            </thead>
            <tbody className="w-full">
              {notifications.map((e: Notification, idx: number) => {
                return (
                  <tr
                    className="h-20 text-lg leading-none text-gray-800 bg-white hover:bg-gray-100 border-b border-t border-gray-100"
                    key={idx}
                  >
                    <td className="pl-16" key={idx}>
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="h-6 w-6"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                      </svg>
                    </td>
                    <td className="pl-16">
                      <p className="font-medium">
                        <span>{e.title}</span>
                      </p>
                    </td>
                    <td className="pl-16">
                      <p className="font-medium">{e.text}</p>
                    </td>
                    <td className="pl-16">
                      <p className="font-medium">{e.date}</p>
                    </td>
                    <td className="pl-16">
                      <Button
                        color="indigo"
                        text="View"
                        loadingText="Opening"
                        onClick={() => viewNotification(e.id!)}
                        loading={loadingNotification}
                      />
                    </td>
                    <td className="pl-1">
                      <Button
                        color="red"
                        text="Delete"
                        loadingText="Deleting"
                        onClick={() => deleteNotification(e.id!)}
                        loading={deleting}
                      />
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default Notifications;
