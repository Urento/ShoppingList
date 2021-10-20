import { useState, useEffect } from "react";
import { Participant } from "../../../types/Participant";
import { API_URL } from "../../../util/constants";

interface LoadInvitationsResponse {
  message: string;
  code: number;
  data: Participant[];
}

export const useLoadInvitations = (condition: boolean) => {
  const [invitations, setInvitations] = useState<Participant[]>([]);
  const [loadingInvitations, setLoadingInvitations] = useState<boolean>(false);

  const loadInvitations = async () => {
    setLoadingInvitations(true);
    const response = await fetch(`${API_URL}/participant/requests`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    });
    const fJson: LoadInvitationsResponse = await response.json();
    setInvitations(fJson.data);
    setLoadingInvitations(false);
  };

  useEffect(() => {
    loadInvitations();
  }, []);

  useEffect(() => {
    if (condition) loadInvitations();
  }, [condition]);

  return {
    invitations,
    setInvitations,
    loadingInvitations,
    setLoadingInvitations,
    loadInvitations,
  };
};
