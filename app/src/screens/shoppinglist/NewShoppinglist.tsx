import clsx from "clsx";
import React from "react";
import { useEffect } from "react";
import { useState } from "react";
import { useHistory } from "react-router-dom";
import swal from "sweetalert";
import { queryClient } from "../..";
import { Sidebar } from "../../components/Sidebar";
import { API_URL } from "../../util/constants";

const regexEmail = /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/;

interface ParticipantsListProps {
  participants: string[];
  removeParticipant(key: number): void;
}

export const ParticipantsList: React.FC<ParticipantsListProps> = ({
  participants,
  removeParticipant,
}) => {
  useEffect(() => {
    console.log(participants);
  }, [participants]);

  if (participants.length > 0) {
    return (
      <div className="flex flex-col">
        <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
            <div className="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th
                      scope="col"
                      className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                    >
                      Name
                    </th>
                    <th
                      scope="col"
                      className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                    >
                      Status
                    </th>
                    <th scope="col" className="relative px-6 py-3">
                      <span className="sr-only">Delete</span>
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="h-6 w-6 text-red-500"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                        />
                      </svg>
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {participants.map((email, key) => (
                    <tr key={key}>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="flex items-center">
                          <div className="ml-4">
                            <div className="text-sm font-medium text-gray-900">
                              {email}
                            </div>
                          </div>
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span
                          className={clsx(
                            `px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                              email != null
                                ? "bg-green-100 text-green-800"
                                : "bg-red-100 text-red-800"
                            }`
                          )}
                        >
                          {email != null ? "Invite Sent" : "Deleted"}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                        <button
                          onClick={() => removeParticipant(key)}
                          className="text-red-600 hover:text-red-900"
                        >
                          Delete
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    );
  } else {
    return <div></div>;
  }
};

interface CreateResponseData {
  message: string;
  success: "true" | "false";
  owner: string;
  position: number;
  participants: string[];
}

interface CreateResponse {
  message: string;
  data: CreateResponseData;
  code: number;
}

export const NewShoppinglist: React.FC = ({}) => {
  const [title, setTitle] = useState("");
  const [participants, setParticipants] = useState<string[]>([]);
  const [currParticipant, setCurrParticipant] = useState("");
  const [loading, setLoading] = useState(false);
  const history = useHistory();

  const handleTitleChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setTitle(e.target.value);
  const handleCurrParticipantChange = (
    e: React.ChangeEvent<HTMLInputElement>
  ) => setCurrParticipant(e.target.value);

  const addParticipant = (email: string) => {
    if (email === "" || participants.indexOf(email) !== -1) return;

    if (!regexEmail.test(email)) {
      swal({
        icon: "error",
        title: "Error inviting Participant",
        text: "Email is not valid",
      });
      return;
    }

    console.log(email);
    if (email.length > 0) {
      setParticipants([...participants, email]);
      setCurrParticipant("");
      swal({
        icon: "success",
        title: "Successfully invited",
        text: "Successfully invited Participant",
      });
    }
  };

  const createList = async (e: any) => {
    e.preventDefault();

    if (title.length <= 0)
      return swal({
        icon: "error",
        title: "Error while creating Shopinglist",
        text: "Title has to be longer than 0 Characters",
      });

    setLoading(true);

    //TODO: remove potential null and undefined entries in array

    const response = await fetch(`${API_URL}/list`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
      body: JSON.stringify({
        title: title,
        participants: participants,
        position: 0,
      }),
    });

    const fJson: CreateResponse = await response.json();
    if (fJson.code != 200)
      return swal({
        icon: "error",
        title: "Error creating Shoppnglist",
        text: fJson.message,
      });

    queryClient.invalidateQueries("shoppinglists");
    swal({
      icon: "success",
      title: "Successfully created",
      text: "Shoppinglist successfully created",
    });
    setLoading(false);
    //TODO: clear timeout?
    setTimeout(() => history.push("/dashboard"), 2000);
  };

  const removeParticipant = (key: number) => {
    delete participants[key];
    setParticipants([...participants]);
    console.log(participants);
  };

  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
      <div className="container mx-auto">
        <div className="flex justify-center px-6 my-12">
          <div className="w-full xl:w-3/4 lg:w-11/12 flex">
            <div className="w-full lg:w-7/12 bg-white p-5 rounded-lg lg:rounded-l-none">
              <form
                className="px-8 pt-6 pb-8 mb-4 bg-white rounded"
                onSubmit={createList}
              >
                <div className="mb-4 md:flex md:justify-between">
                  <div className="mb-4 md:mr-2 md:mb-0">
                    <label
                      className="mb-2 text-sm font-bold text-gray-700"
                      htmlFor="title"
                    >
                      Title
                    </label>
                    <input
                      className="w-full px-3 py-2 text-sm leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline"
                      id="title"
                      type="text"
                      placeholder="Title"
                      onChange={handleTitleChange}
                    />
                  </div>
                  <div className="md:ml-2">
                    <label
                      className="mb-2 text-sm font-bold text-gray-700"
                      htmlFor="username"
                    >
                      Participants
                    </label>
                    <input
                      className="w-full px-3 py-2 text-sm leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline"
                      id="username"
                      type="text"
                      placeholder="Username"
                      onChange={handleCurrParticipantChange}
                    />
                    <button
                      className="w-full px-4 py-2 font-bold text-white bg-indigo-500 rounded-full hover:bg-indigo-700 focus:outline-none focus:shadow-outline"
                      type="button"
                      onClick={() => addParticipant(currParticipant)}
                    >
                      Invite Participant
                    </button>
                  </div>
                </div>
                <div className="mb-6 text-center">
                  <button
                    className="w-full px-4 py-2 font-bold text-white bg-green-500 rounded-full hover:bg-green-700 focus:outline-none focus:shadow-outline"
                    type="submit"
                    onSubmit={createList}
                  >
                    {loading ? (
                      <svg
                        className="loading-svg justify-center flex"
                        viewBox="25 25 50 50"
                      >
                        <circle
                          className="loading-circle"
                          cx="50"
                          cy="50"
                          r="20"
                        ></circle>
                      </svg>
                    ) : (
                      "Create"
                    )}
                  </button>
                </div>
              </form>
              <ParticipantsList
                participants={participants}
                removeParticipant={removeParticipant}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
