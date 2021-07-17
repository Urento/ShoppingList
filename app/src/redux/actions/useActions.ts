import { StackNavigationProp } from "@react-navigation/stack";
import { Dispatch } from "redux";
import {
  API_URL,
  AUTH_API_URL,
  HomeStackParametersList,
} from "../../util/constants";
import { ACTION_TYPES } from "../types";

export type UserData = {
  email: string;
  password: string;
  token: string;
  loggedIn: boolean;
};

export const loginUser =
  async (
    userData: UserData,
    navigation: StackNavigationProp<HomeStackParametersList>
  ) =>
  async (dispatch: Dispatch) => {
    dispatch({ type: ACTION_TYPES.LOADING_UI });

    try {
      const f = await fetch(AUTH_API_URL, {
        method: "POST",
        body: JSON.stringify(userData),
      });
      const fJson = await f.json();
      const token = `Bearer ${fJson.token}`;
      localStorage.setItem("token", token);
      //TODO: set authorization header in every request

      const userD = await getUserData();
      dispatch({ type: ACTION_TYPES.SET_USER, payload: userD });
      dispatch({ type: ACTION_TYPES.CLEAR_ERRORS });
      navigation.navigate("Home");
    } catch (err) {
      //report err to sentry
      console.log(err);
      dispatch({ type: ACTION_TYPES.SET_ERRORS, payload: err });
    }
  };

export const getUserData = () => async (dispatch: any) => {
  dispatch({ type: ACTION_TYPES.LOADING_USER });

  try {
    const token = localStorage.getItem("token");
    const f = await fetch(`${API_URL}/user`, {
      method: "GET",
      headers: {
        Authentication: token!,
        "Content-Type": "application/json",
      },
    });
    const fJson = await f.json();
    dispatch({ type: ACTION_TYPES.SET_USER, payload: fJson.data });
    return fJson;
  } catch (err) {
    console.log(err);
    return err;
  }
};

export const logout =
  (navigation: StackNavigationProp<HomeStackParametersList>) =>
  async (dispatch: any) => {
    const token = localStorage.getItem("token");
    try {
      fetch(`${API_URL}/logout`, {
        method: "GET",
        headers: {
          Authorization: token!,
          "Content-Type": "application/json",
        },
      });
    } catch (err) {
      console.log(err);
    }
    localStorage.removeItem("token");
    dispatch({ type: ACTION_TYPES.SET_UNAUTHENTICATED });
    navigation.navigate("Login");
  };
