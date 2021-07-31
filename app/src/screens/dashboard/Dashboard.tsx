import React, { useState } from "react";
import { Loading } from "../../components/Loading";
import useAuthCheck from "../../hooks/useAuthCheck";

export const Dashboard: React.FC = () => {
  //const [redirect, setRedirect] = useState(false);

  const authStatus = useAuthCheck();
  console.log(authStatus);
  if (authStatus === "fail") {
    localStorage.removeItem("authenticated");
    //TODO: actually redirect back to login
  }

  return <>{authStatus === "pending" && <Loading />}Dashboard</>;
};
