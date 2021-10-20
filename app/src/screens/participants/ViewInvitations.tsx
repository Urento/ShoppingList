import { useHistory } from "react-router";
import { useState } from "react";
import { Sidebar } from "../../components/Sidebar";
import { useLoadInvitations } from "./hooks/useLoadInvitations";
import { Button } from "../../components/Button";
import { Participant } from "../../types/Participant";
import { API_URL } from "../../util/constants";

const ViewInvitations: React.FC = () => {
  const history = useHistory();
  const [refetch, setRefetch] = useState<boolean>(false);
  const {
    invitations,
    setInvitations,
    loadingInvitations,
    setLoadingInvitations,
    loadInvitations,
  } = useLoadInvitations(refetch);

  const unixToDate = (timestamp: number) => {
    const a = new Date(timestamp);
    var months = [
      "Jan",
      "Feb",
      "Mar",
      "Apr",
      "May",
      "Jun",
      "Jul",
      "Aug",
      "Sep",
      "Oct",
      "Nov",
      "Dec",
    ];
    var year = a.getFullYear();
    var month = months[a.getMonth()];
    var date = a.getDate();
    var hour = a.getHours();
    var min = a.getMinutes() < 9 ? "0" + a.getMinutes() : a.getMinutes();
    return date + " " + month + " " + year + " " + hour + ":" + min;
  };

  const acceptInvitation = async (id: number) => {
    await fetch(`${API_URL}/participant/requests`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify({
        id: id,
      }),
      credentials: "include",
    });
    loadInvitations();
  };

  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
      <div className="container mx-auto py-10 md:w-4/5 w-11/12">
        <div className="bg-white px-4 md:px-10 pt-5 md:pt-7 pb-5 overflow-y-auto">
          {invitations.length > 0 && (
            <Button
              text="Deny all"
              loadingText="Marking all notifications as read..."
              onClick={() => console.log("")}
              className="inline-flex sm:ml-3 mt-4 sm:mt-0 items-start justify-start px-6 py-3 bg-red-600 hover:bg-red-500 text-white focus:outline-none rounded"
              loading={false}
            />
          )}
          {invitations.length <= 0 && <p>You have no invitations!</p>}
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
              {invitations.map((e: Participant, idx: number) => {
                return (
                  <tr
                    className="h-20 text-lg leading-none text-gray-800 bg-white hover:bg-gray-100 border-b border-t border-gray-100"
                    key={idx}
                  >
                    <td className="pl-16">
                      <p className="font-medium">
                        <span className="font-bold">From:</span>{" "}
                        {e.request_from}
                      </p>
                    </td>
                    <td className="pl-16">
                      <p className="font-medium">{unixToDate(e.created_on!)}</p>
                    </td>
                    <td className="pl-1">
                      <Button
                        color="green"
                        text="Accept"
                        loadingText="Accepting..."
                        onClick={() => acceptInvitation(e.id!)}
                        loading={false}
                      />
                    </td>
                    <td className="pl-1">
                      <Button
                        color="red"
                        text="Deny"
                        loadingText="Denying..."
                        onClick={() => console.log(e.id!)}
                        loading={false}
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

export default ViewInvitations;
