import { useState } from "react";
import { useEffect } from "react";
import { getUser, UserData } from "../storage/UserStorage";

export const useGetUserData = () => {
  const [user, setUser] = useState<UserData>();

  useEffect(() => {
    const fetch = async () => {
      const userObj = await getUser();
      setUser(userObj);
    };

    fetch();
  }, []);

  return user;
};
