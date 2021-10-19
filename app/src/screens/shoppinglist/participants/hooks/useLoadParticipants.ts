import { useState } from "react";
import { useEffect } from "react";
import { Participant, RequestsFromList } from "../../../../types/Participant";
import { API_URL } from "../../../../util/constants";

export const useLoadParticipants = (id: number, condition: boolean) => {
  const [participants, setParticipants] = useState<Participant[]>([]);
  const [loadingParticipants, setLoadingParticipants] = useState<boolean>(true);

  const loadParticipants = async () => {
    const response = await fetch(`${API_URL}/participants/${id}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    });
    const fJson: RequestsFromList = await response.json();
    if (fJson.code === 200) setParticipants(fJson.data);
    setLoadingParticipants(false);
  };

  useEffect(() => {
    loadParticipants();
  }, []);

  useEffect(() => {
    if (condition) loadParticipants();
  }, [condition]);

  return {
    participants,
    setParticipants,
    loadingParticipants,
    setLoadingParticipants,
  };
};
