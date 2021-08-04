import React from "react";
import { useParams } from "react-router";
import { Sidebar } from "../../components/Sidebar";

interface Params {
  id: string;
}

export const ViewShoppinglist: React.FC = ({}) => {
  const { id } = useParams<Params>();

  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
    </div>
  );
};
