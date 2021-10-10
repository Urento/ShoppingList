import { useHistory } from "react-router";
import { useState } from "react";
import { Button } from "../../components/Button";
import { Loading } from "../../components/Loading";
import { Sidebar } from "../../components/Sidebar";
import useAuthCheck from "../../hooks/useAuthCheck";

const Notifications: React.FC = ({}) => {
  const history = useHistory();
  const authStatus = useAuthCheck();
  const [loading, setLoading] = useState<boolean>(false);

  if (authStatus === "fail") {
    localStorage.removeItem("authenticated");
    history.push("/");
  }

  if (authStatus === "pending") {
    return <Loading withSidebar />;
  }

  const deleteNotification = async (id: number) => {
    setLoading(true);
  };

  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
      <div className="container mx-auto py-10 md:w-4/5 w-11/12">
        <div className="bg-white px-4 md:px-10 pt-4 md:pt-7 pb-5 overflow-y-auto">
          <table className="w-full whitespace-nowrap">
            <thead>
              <tr className="h-16 w-full text-sm leading-none ">
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
              </tr>
            </thead>
            <tbody className="w-full">
              <tr className="h-20 text-lg leading-none text-gray-800 bg-white hover:bg-gray-100 border-b border-t border-gray-100">
                <td className="pl-48">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-6 w-6 text-green-600"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                </td>
                <td className="pl-48">
                  <p className="font-medium">
                    <span className="text-red-600">fghdfghj</span>
                  </p>
                </td>
                <td className="pl-48">
                  <p className="font-medium">dfgh</p>
                </td>
                <td className="pl-48">
                  <Button
                    color="red"
                    text="Delete"
                    loadingText="Deleting"
                    onClick={() => deleteNotification(123)}
                    loading={loading}
                  />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default Notifications;
