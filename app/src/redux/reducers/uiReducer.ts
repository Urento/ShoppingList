import { ACTION_TYPES } from "../types";

const initialState = {
  loading: false,
  errors: null,
};

export default function (state = initialState, action: any) {
  switch (action.type) {
    case ACTION_TYPES.SET_ERRORS:
      return {
        ...state,
        loading: false,
        errors: action.payload,
      };
    case ACTION_TYPES.CLEAR_ERRORS:
      return {
        ...state,
        loading: false,
        errors: null,
      };
    case ACTION_TYPES.LOADING_UI:
      return {
        ...state,
        loading: true,
      };
    default:
      return state;
  }
}
