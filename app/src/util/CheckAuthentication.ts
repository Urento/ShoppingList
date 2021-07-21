import jwtDecode from "jwt-decode";
import { getUserData, logout } from "../redux/actions/useActions";
import { ACTION_TYPES } from "../redux/types";
import store from "../redux/store";
import { getJWTToken } from "src/realm/realm";

export const CheckAuthentication = () => {
  const authToken = getJWTToken();

  if (authToken) {
    const decodedToken: any = jwtDecode(authToken);
    if (decodedToken.exp * 1000 < Date.now()) {
      store.dispatch(logout());
    } else {
      store.dispatch({ type: ACTION_TYPES.SET_AUTHENTICATED });
      store.dispatch(getUserData());
    }
  }
};
