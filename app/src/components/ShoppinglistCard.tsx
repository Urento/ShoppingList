import React, { useState } from "react";
import { useEffect } from "react";
import { useQuery } from "react-query";
import swal from "sweetalert";
import { API_URL } from "../util/constants";
import { Button } from "./Button";

interface ListData {
  title: string;
  items: string[];
  owner: string;
  participants: string[];
  position: number;
  created_on: string;
  modified_on: number;
}

export const ShoppinglistCard: React.FC = ({}) => {
  const [shoppinglists, setShoppinglists] = useState<any>();

  const { isLoading, error, data } = useQuery<any, Error>("shoppinglists", () =>
    fetch(`${API_URL}/lists`, {
      method: "GET",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
    })
      .then((res: Response) => res.json())
      .then((data) => setShoppinglists(data.data))
  );

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

  return (
    <div>
      {shoppinglists.map((e: ListData) => (
        <div className="max-w-md py-4 px-8 bg-gray-800 shadow-lg rounded-lg">
          <div>
            <h2 className="text-white text-3xl font-semibold">{e.title}</h2>
            <p className="mt-2 text-white">
              <span className="font-bold">Participants</span>: <br />
              {e.participants
                .toString()
                .split(" ")
                .join(",")
                .replace(",", ", ")
                .substring(0, 100)}
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
            <a href="#" className="text-xl font-medium text-white">
              View Shoppinglist
            </a>
          </div>
        </div>
      ))}
    </div>
  );
};
