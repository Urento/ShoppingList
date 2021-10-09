import React from "react";
import { useQuery } from "react-query";
import { useParams } from "react-router";
import { useHistory } from "react-router-dom";
import { ItemsTable } from "../../components/ItemsTable";
import { Sidebar } from "../../components/Sidebar";
import { ListResponse, ListResponseData } from "../../types/Shoppinglist";
import { API_URL } from "../../util/constants";

interface Params {
  id: string;
}

const ViewShoppinglist: React.FC = ({}) => {
  const { id } = useParams<Params>();

  const { isLoading, data, error } = useQuery<ListResponseData, Error>(
    `shoppinglist_${id}`,
    () => fetchData(id),
    { refetchOnWindowFocus: false }
  );

  if (error) {
    return <div>Error while loading data: {error}</div>;
  }

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (data === null) {
    return <div>Error while loading data... Try again!</div>;
  }

  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
      <div className="container mx-auto">
        <div className="flex justify-center px-6 my-12">
          <div className="w-full xl:w-3/4 lg:w-11/12 flex">
            <ItemsTable data={data!} />
          </div>
        </div>
      </div>
    </div>
  );
};

const fetchData = async (id: string) => {
  const response = await fetch(`${API_URL}/list/${id}`, {
    method: "GET",
    headers: { "Content-Type": "application/json", Accept: "application/json" },
    credentials: "include",
  });
  const fJson: ListResponse = await response.json();
  return fJson.data;
};

export default ViewShoppinglist;
