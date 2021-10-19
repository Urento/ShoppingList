import { useHistory, useParams } from "react-router";
import { useState } from "react";
import { Button } from "../../../components/Button";
import { Loading } from "../../../components/Loading";
import { Sidebar } from "../../../components/Sidebar";
import {
  AddParticipantResponse,
  Participant,
} from "../../../types/Participant";
import { useLoadParticipants } from "./hooks/useLoadParticipants";
import swal from "sweetalert";
import { API_URL } from "../../../util/constants";

export interface Params {
  id: string;
}

const Participants: React.FC = () => {
  const { id } = useParams<Params>();
  const history = useHistory();
  const [refresh, setRefresh] = useState<boolean>(false);
  const [invitingParticipant, setInvitingParticipant] =
    useState<boolean>(false);
  const [removingParticipant, setRemovingParticipant] = useState({
    id: 0,
    loading: false,
  });
  const {
    participants,
    setParticipants,
    loadingParticipants,
    setLoadingParticipants,
  } = useLoadParticipants(parseInt(id), refresh);

  if (loadingParticipants) return <Loading withSidebar />;

  const validateEmail = (email: string) => {
    const re =
      /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(String(email).toLowerCase());
  };

  const inviteNewParticipant = async () => {
    setInvitingParticipant(true);

    await swal({
      text: "Create a new Participant", //@ts-ignore
      content: "input",
      buttons: ["Cancel", "Invite"],
      closeOnEsc: false,
      closeOnClickOutside: false,
    }).then(async (email: string) => {
      if (email === "" || !email) return;
      if (!validateEmail(email)) {
        swal.close!();
        swal({
          icon: "error",
          title: "Email is not valid",
        });
        setInvitingParticipant(false);
        return;
      }

      const response = await fetch(`${API_URL}/participant`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: JSON.stringify({
          parentListId: parseInt(id),
          email: email,
        }),
        credentials: "include",
      });
      const fJson: AddParticipantResponse = await response.json();
      if (
        fJson.data &&
        fJson.data.error === "participant is already included"
      ) {
        swal.close!();
        swal({
          icon: "error",
          title: "Participant is already included",
        });
        setInvitingParticipant(false);
        return;
      }
      setRefresh(true);
      setInvitingParticipant(false);
      setTimeout(() => setRefresh(false), 1000);
    });
  };

  const deleteParticipant = async (participantId: number) => {
    setRemovingParticipant({
      id: participantId,
      loading: true,
    });
    await fetch(`${API_URL}/participant/${id}/${participantId}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    });
    setRefresh(true);
    setRemovingParticipant({ id: 0, loading: false });
    setTimeout(() => setRefresh(false), 1000);
  };

  const isRemoving = (id: number): boolean => {
    return removingParticipant.id === id && removingParticipant.loading;
  };

  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
      <div className="container mx-auto py-10 md:w-4/5 w-11/12">
        <h1 className="text-2xl text-center">Participants</h1>
        <div className="bg-white px-4 md:px-10 pt-5 md:pt-7 pb-5 overflow-y-auto">
          <br />
          <Button
            text="Go Back"
            onClick={history.goBack}
            className="inline-flex sm:ml-3 mt-4 sm:mt-0 items-start justify-start px-6 py-3 bg-indigo-600 hover:bg-indigo-500 text-white focus:outline-none rounded"
          />
          <Button
            text="Invite new Participant"
            loadingText="Inviting new Participant..."
            onClick={inviteNewParticipant}
            className="inline-flex sm:ml-3 mt-4 sm:mt-0 items-start justify-start px-6 py-3 bg-green-600 hover:bg-green-500 text-white focus:outline-none rounded"
            loading={invitingParticipant}
          />
          {participants.length <= 0 && <p>No Participants</p>}
          <table className="w-full whitespace-nowrap">
            <thead>
              <tr className="h-16 w-full text-sm leading-none ">
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
              </tr>
            </thead>
            <tbody className="w-full">
              {participants.map((e: Participant, idx: number) => {
                const removing = isRemoving(e.id!);
                return (
                  <tr
                    className="h-20 text-lg leading-none text-gray-800 bg-white hover:bg-gray-100 border-b border-t border-gray-100"
                    key={idx}
                  >
                    <td className="pl-16">
                      <p className="font-medium">{e.email}</p>
                    </td>
                    <td className="pl-16">
                      <p className="font-medium">
                        {e.status === "accepted" ? (
                          <span className="text-green-500">Accepted</span>
                        ) : (
                          <span className="text-gray-500">Pending</span>
                        )}
                      </p>
                    </td>
                    <td className="pl-16">
                      <Button
                        text="Remove"
                        loadingText="Removing..."
                        loading={removing}
                        color="red"
                        onClick={() => deleteParticipant(e.id!)}
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

export default Participants;
