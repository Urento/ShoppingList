import React from "react";
import Login from "./components/form/Login";

export type Props = {
  loggedIn: boolean;
};

const login = (e: any) => {
  e.preventDefault();
};

const onChangeEmail = (e: any) => {};

const App: React.FC<Props> = ({ loggedIn = false }) => {
  return <Login loggedIn={loggedIn} />;
};
export default App;
