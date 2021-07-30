import React, { useEffect, useState } from "react";
import { Redirect } from "react-router-dom";
import useAuthCheck from "../../hooks/useAuthCheck";

export const Dashboard: React.FC = () => {
  const [redirect, setRedirect] = useState(false);

  const authStatus = useAuthCheck();
  if (authStatus != "success") {
    localStorage.removeItem("authenticated");
    setRedirect(true);
  }

  return (
    <div>
      {redirect && <Redirect to="/"></Redirect>}
      <h1>dfgsfdgasdfsadf</h1>
    </div>
  );
};
