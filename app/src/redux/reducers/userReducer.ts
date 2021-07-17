import { ACTION_TYPES } from "../types";

const initialState = {
  authenticated: false,
  credentials: {},
  loading: false,
};

export default function (state = initialState, action: any) {
  switch (action.type) {
    case ACTION_TYPES.SET_AUTHENTICATED:
      return {
        ...state,
        authenticated: true,
      };
    case ACTION_TYPES.SET_UNAUTHENTICATED:
      return {
        ...state,
        authenticated: false,
        loading: false,
        ...action.payload,
      };
    case ACTION_TYPES.LOADING_USER:
      return {
        ...state,
        loading: true,
      };
    default:
      return state;
  }
}
