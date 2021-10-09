import { API_URL } from "../util/constants";

export type UserData = {
  created_on: number;
  modified_on: number;
  deleted_at: number | null;
  email: string;
  username: string;
  twoFactorAuthentication: boolean;
  rank: string;
  email_verified: boolean;
};

interface UserDataResponseData {
  created_on: number;
  modified_on: number;
  deleted_at: number | null;
  id: number;
  e_mail: string;
  email_verified: boolean;
  username: string;
  rank: string;
  two_factor_authentication: boolean;
}

interface UserDataResponse {
  message: string;
  code: number;
  data: UserDataResponseData;
}

export const saveUser = (data: UserData) => {
  localStorage.removeItem("user");
  localStorage.setItem("user", JSON.stringify(data));
};

export const getUser = async (): Promise<UserData> => {
  if (
    localStorage.getItem("user") === null ||
    localStorage.getItem("user") === undefined
  ) {
    const newUserObj = await fetchUser();
    return newUserObj;
  }

  const userObj: UserData = JSON.parse(localStorage.getItem("user")!);
  const userDataObj: UserData = {
    email: userObj.email,
    username: userObj.username,
    twoFactorAuthentication: userObj.twoFactorAuthentication,
    created_on: userObj.created_on,
    deleted_at: userObj.deleted_at,
    email_verified: userObj.email_verified,
    modified_on: userObj.modified_on,
    rank: userObj.rank,
  };

  return userDataObj;
};

export const updateTwoFA = async (status: boolean) => {
  if (localStorage.getItem("user") === null) {
    await fetchUser();
    return;
  }

  const prevUserObj: UserData = JSON.parse(localStorage.getItem("user")!);
  localStorage.removeItem("user");
  const userObj: UserData = {
    email: prevUserObj.email,
    username: prevUserObj.username,
    twoFactorAuthentication: status,
    created_on: prevUserObj.created_on,
    deleted_at: prevUserObj.deleted_at,
    email_verified: prevUserObj.email_verified,
    modified_on: prevUserObj.modified_on,
    rank: prevUserObj.rank,
  };
  localStorage.setItem("user", JSON.stringify(userObj));
};

export const removeUser = () => {
  localStorage.removeItem("user");
};

export const fetchUser = async () => {
  const response = await fetch(`${API_URL}/auth/user`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    credentials: "include",
  });
  const fJson: UserDataResponse = await response.json();
  if (fJson.code !== 200) throw new Error("Error while fetching user data");

  const newUserObj: UserData = {
    email: fJson.data.e_mail,
    email_verified: fJson.data.email_verified,
    created_on: fJson.data.created_on,
    deleted_at: fJson.data.deleted_at,
    modified_on: fJson.data.modified_on,
    rank: fJson.data.rank,
    twoFactorAuthentication: fJson.data.two_factor_authentication,
    username: fJson.data.username,
  };
  localStorage.setItem("user", JSON.stringify(newUserObj));
  return newUserObj;
};

export const isLoggedIn = () => {
  return Boolean(localStorage.getItem("authenticated"));
};
