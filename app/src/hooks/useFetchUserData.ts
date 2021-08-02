import { useEffect, useState } from "react";
import { API_URL } from "../util/constants";

export const useFetchUserData = (conditionToRefetch: boolean | undefined) => {
  const [user, setUser] = useState<any>();

  useEffect(() => {
    const fetchUser = async () => {
      const response = await fetch(`${API_URL}/auth/user`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        credentials: "include",
      });
      const fJson = await response.json();
      setUser(fJson.data);
    };

    fetchUser();
  }, [conditionToRefetch]);

  return user;
};
