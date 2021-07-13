import React from "react";
import Login from "./components/form/Login";

export type Props = {
  loggedIn: boolean;
};

const App: React.FC<Props> = ({ loggedIn = false }) => {
  return <Login loggedIn={loggedIn} />;
};
export default App;
