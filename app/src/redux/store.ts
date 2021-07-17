import { applyMiddleware, combineReducers, compose, createStore } from "redux";
import thunk from "redux-thunk";
import uiReducer from "./reducers/uiReducer";
import userReducer from "./reducers/userReducer";

const initialState = {};
const middleware = [thunk];

//for redux devtool
declare global {
  interface Window {
    __REDUX_DEVTOOLS_EXTENSION__?: typeof compose;
  }
}

const reducer = combineReducers({
  user: userReducer,
  UI: uiReducer,
});

const store = createStore(
  reducer,
  initialState,
  compose(
    applyMiddleware(...middleware),
    (window.__REDUX_DEVTOOLS_EXTENSION__ &&
      window.__REDUX_DEVTOOLS_EXTENSION__()) as any
  )
);

export default store;
