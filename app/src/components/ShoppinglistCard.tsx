import React, { useState } from "react";
import { useEffect } from "react";
import { useQuery } from "react-query";
import swal from "sweetalert";
import { API_URL } from "../util/constants";

interface ListData {
  title: string;
  items: string[];
  owner: string;
  participants: string[];
  position: number;
}

interface Props {}

export const ShoppinglistCard: React.FC<Props> = ({}) => {
  const { isLoading, error, data } = useQuery<ListData[], Error>(
    "shoppinglists",
    () =>
      fetch(`${API_URL}/lists`, {
        method: "GET",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
      }).then((res: any) => res.json())
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

  console.log(data);

  return (
    <div>
      <div className="max-w-screen md:w-3/4 mx-auto">
        <div className="flex flex-row space-y-2 items-center justify-center h-full py-4 bg-gray-800 rounded-xl space-x-10">
          <div className="w-2/3">
            <p className="w-full text-2xl font-semibold text-white">
              We love pixels
            </p>
            <p className="w-full pb-8 text-sm tracking-wide leading-tight text-white">
              The card layouts can vary to support the types of content they
              contain.
            </p>
            <div className="rounded w-1/3">
              <div className="opacity-95 border rounded-lg border-white px-4">
                <p className="m-auto inset-0 text-sm font-medium leading-normal text-center text-white py-2">
                  License
                </p>
              </div>
            </div>
          </div>
          <div className="w-auto h-">
            <img
              className="flex-1 h-full rounded-lg"
              src="https://via.placeholder.com/96x136"
            />
          </div>
        </div>
      </div>
    </div>
  );
};
