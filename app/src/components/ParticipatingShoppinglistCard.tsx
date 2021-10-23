import { useState } from "react";
import { useQuery } from "react-query";
import { useHistory } from "react-router";
import swal from "sweetalert";
import { Shoppinglist } from "../types/Shoppinglist";
import { API_URL } from "../util/constants";
import { Loading } from "./Loading";
import { DeleteResponse, NoItemsToDisplay } from "./ShoppinglistCard";

export const ParticipatingShoppinglistCard: React.FC = () => {
  const [participatingShoppinglists, setParticipatingShoppinglists] = useState<
    Shoppinglist[]
  >([]);
  const history = useHistory();

  const { isLoading, error, isFetching, refetch } = useQuery<any, Error>(
    "shoppinglists_by_participation",
    async () =>
      await fetch(`${API_URL}/listsByParticipation`, {
        method: "GET",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
      })
        .then((res: Response) => res.json())
        .then((data) => setParticipatingShoppinglists(data.data)),
    { refetchOnWindowFocus: false }
  );

  if (isFetching) return <Loading />;
  if (isLoading) return <Loading />;

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

    refetch();

    return swal({
      icon: "success",
      title: "Successfully deleted",
      text: "Successfully deleted the Shoppinglist",
    });
  };

  const leaveShoppinglist = async (id: number) => {
    const alert = await swal({
      icon: "warning",
      title: "Are you sure you want to leave the shoppinglist?",
      buttons: ["No, don't leave!", "Yes, leave!"],
    });

    if (alert) {
      const response = await fetch(`${API_URL}/participant/list/leave`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: JSON.stringify({ id: id }),
        credentials: "include",
      });

      if (response) refetch();
    }
  };

  return (
    <div className="flex flex-wrap">
      {participatingShoppinglists.length <= 0 && (
        <p>You are not participating in any other Shoppinglists</p>
      )}
      {participatingShoppinglists.length > 0 &&
        participatingShoppinglists.map((e: Shoppinglist) => (
          <div className="pt-2 pl-2">
            <div className="max-w-md py-4 px-8 bg-gray-800 shadow-lg rounded-lg">
              <div className="justify-center md:justify-end -m-3.5 pl-96">
                <button onClick={() => leaveShoppinglist(e.id)}>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-9 object-cover rounded-full text-red-600"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M12 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2M3 12l6.414 6.414a2 2 0 001.414.586H19a2 2 0 002-2V7a2 2 0 00-2-2h-8.172a2 2 0 00-1.414.586L3 12z"
                    />
                  </svg>
                </button>
              </div>
              <div>
                <h2 className="text-white text-3xl font-semibold">{e.title}</h2>
                <p className="mt-2 text-white">
                  <span className="font-bold">Participants</span>:{" "}
                  {e.participants ? e.participants.length + 1 : 1}
                </p>
                <p className="mt-2 text-white">
                  <span className="font-bold">Last Edited</span>:{" "}
                  {e.modified_on && unixToDate(e.modified_on!)}
                </p>
                <p className="mt-2 text-white">
                  <span className="font-bold">Created</span>:{" "}
                  {e.created_on && unixToDate(e.created_on!)}
                </p>
              </div>
              <div className="flex justify-end mt-4">
                <button
                  onClick={() => history.push(`/list/${e.id}`)}
                  className="text-lg font-bold text-white"
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
