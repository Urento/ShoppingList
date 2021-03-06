import React from "react";
import { useHistory } from "react-router-dom";
import { Loading } from "../components/Loading";
import { ParticipatingShoppinglistCard } from "../components/ParticipatingShoppinglistCard";
import { ShoppinglistCard } from "../components/ShoppinglistCard";
import { Sidebar } from "../components/Sidebar";
import useAuthCheck from "../hooks/useAuthCheck";

const Dashboard: React.FC = () => {
  const history = useHistory();
  const authStatus = useAuthCheck();

  if (authStatus === "fail") {
    localStorage.removeItem("authenticated");
    history.push("/");
  }

  if (authStatus === "pending") return <Loading withSidebar />;

  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
      <div className="container mx-auto py-10 md:w-4/5 w-11/12">
        <div className="w-full h-full rounded">
          <div className="flex items-center justify-between">
            <ShoppinglistCard />
          </div>
          <br />
          <hr />
          <br />
          <div className="flex items-center justify-between">
            <ParticipatingShoppinglistCard />
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
