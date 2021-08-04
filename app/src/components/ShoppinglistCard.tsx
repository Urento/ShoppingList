import React, { useState } from "react";
import { useMutation, useQuery } from "react-query";
import { useHistory } from "react-router-dom";
import swal from "sweetalert";
import { queryClient } from "..";
import { API_URL } from "../util/constants";

interface ListData {
  id: number;
  title: string;
  items: string[];
  owner: string;
  participants: string[];
  position: number;
  created_on: string;
  modified_on: number;
}

interface DeleteResponseData {
  success: "true" | "false";
  message: string;
}

interface DeleteResponse {
  data: DeleteResponseData;
  code: number;
  message: string;
}

export const ShoppinglistCard: React.FC = ({}) => {
  const [shoppinglists, setShoppinglists] = useState<any>();
  const history = useHistory();

  const { isLoading, error, isFetching } = useQuery<any, Error>(
    "shoppinglists",
    () =>
      fetch(`${API_URL}/lists`, {
        method: "GET",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
      })
        .then((res: Response) => res.json())
        .then((data) => setShoppinglists(data.data)),
    { refetchOnWindowFocus: false }
  );

  const addMutation = useMutation(
    () =>
      fetch(`${API_URL}/auth/user`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        credentials: "include",
      }),
    {
      onSuccess: () => queryClient.invalidateQueries("shoppinglists"),
    }
  );

  if (isFetching) {
    return <div>Fetching....</div>;
  }

  if (isLoading) {
    return <div>Loading....</div>;
  }

  if (error) {
    swal({
      icon: "error",
      text: "Error while getting the Shoppinglists",
      title: "Error while getting Shoppinglists",
    });
  }

  const unixToDate = (UNIX_timestamp: number) => {
    const a = new Date(UNIX_timestamp);
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
    var time = date + " " + month + " " + year + " " + hour + ":" + min;
    return time;
  };

  const showDeleteListModal = (id: number) => {
    swal({
      icon: "warning",
      title: "Are you sure?",
      text: "Are you sure you want to delete the Shoppinglist?",
      dangerMode: true,
      buttons: ["No, dont delete!", "Yes, delete!"],
    }).then(async (willDelete: boolean) => {
      if (willDelete) await deleteList(id);
    });
  };

  const deleteList = async (id: number) => {
    const response = await fetch(`${API_URL}/list/${id}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    });

    const fJson: DeleteResponse = await response.json();

    if (fJson.code !== 200)
      return swal({
        icon: "error",
        title: "Error while deleting",
        text: "Error while deleting Shoppinglist. Try again later!",
      });

    addMutation.mutate();

    return swal({
      icon: "success",
      title: "Successfully deleted",
      text: "Successfully deleted the Shoppinglist",
    });
  };

  return (
    <div className="flex flex-wrap">
      {shoppinglists.map((e: ListData) => (
        <div className="pt-2 pl-2">
          <div className="max-w-md py-4 px-8 bg-gray-800 shadow-lg rounded-lg">
            <div className="justify-center md:justify-end -m-3.5 pl-96">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-9 object-cover rounded-full text-red-600"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                onClick={() => showDeleteListModal(e.id)}
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                />
              </svg>
            </div>
            <div>
              <h2 className="text-white text-3xl font-semibold">{e.title}</h2>
              <p className="mt-2 text-white">
                <span className="font-bold">Participants</span>:{" "}
                {e.participants.length}
              </p>
              <p className="mt-2 text-white">
                <span className="font-bold">Last Edited</span>:{" "}
                {unixToDate(e.modified_on)}
              </p>
              <p className="mt-2 text-white">
                <span className="font-bold">Created</span>:{" "}
                {unixToDate(e.modified_on)}
              </p>
            </div>
            <div className="flex justify-end mt-4">
              <button
                onClick={() => history.push(`/list/${e.id}`)}
                className="text-lg font-bold font-medium text-white"
              >
                View Shoppinglist
              </button>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};
