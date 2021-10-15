import React from "react";
import { useQuery } from "react-query";
import { useParams } from "react-router";
import { useState } from "react";
import { Button } from "../../components/Button";
import { Sidebar } from "../../components/Sidebar";
import { Item, ListResponse, ListResponseData } from "../../types/Shoppinglist";
import { API_URL } from "../../util/constants";
import { Loading } from "../../components/Loading";
import { ReactSortable } from "react-sortablejs";
import swal from "sweetalert";

interface Params {
  id: string;
}

interface DeletingItemState {
  id: number;
  loading: boolean;
}

const ViewShoppinglist: React.FC = ({}) => {
  const { id } = useParams<Params>();
  const [creatingItem, setCreatingItem] = useState(false);
  const [deletingList, setDeletingList] = useState(false);
  const [deletingItem, setDeletingItem] = useState<DeletingItemState>({
    id: 0,
    loading: false,
  });
  const [items, setItems] = useState<Item[]>([]);

  const { isLoading, data, error, refetch } = useQuery<ListResponseData, Error>(
    `shoppinglist_${id}`,
    async () => await fetchData(id),
    { refetchOnWindowFocus: false }
  );

  if (error) return <Loading withSidebar />;
  if (isLoading) return <Loading withSidebar />;
  if (!data) return <Loading withSidebar />;
  //if (!isLoading) setItems(data?.items);

  const createItem = async () => {
    setCreatingItem(true);
    let pos: number;
    let position: number;

    console.log(data?.items[data.items.length - 1].position);

    try {
      pos = data?.items[data?.items.length].position;
      position = pos + 1;
    } catch {
      position = 1;
    }

    await swal({
      text: "Create a new Item", //@ts-ignore
      content: "input",
      button: {
        text: "Create",
        closeModal: "Dont Create",
      },
    }).then(async (title: string) => {
      if (title === "" || !title) return;
      await fetch(`${API_URL}/list/items`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: JSON.stringify({
          id: parseInt(id),
          title: title,
          position: position,
        }),
        credentials: "include",
      });
      refetch();
    });
  };

  const deleteList = async () => {
    setDeletingList(true);
  };

  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
      <div className="container mx-auto py-10 md:w-4/5 w-11/12">
        <div className="bg-white px-4 md:px-10 pt-5 md:pt-7 pb-5 overflow-y-auto">
          <h1 className="text-2xl">{data?.title}</h1>
          <h1 className="text-2xl">
            {items.length} {items.length === 1 ? "Item" : "Items"}
          </h1>
          <h1 className="text-2xl">
            {data?.participants.length === 0 ? 1 : data?.participants.length}{" "}
            {data?.participants.length === 1 || data?.participants.length === 0
              ? "Participant"
              : "Participants"}
          </h1>
          <br />
          <Button
            text="Delete Shoppinglist"
            loadingText="Creating new Item..."
            onClick={deleteList}
            className="inline-flex sm:ml-3 mt-4 sm:mt-0 items-start justify-start px-6 py-3 bg-red-600 hover:bg-red-500 text-white focus:outline-none rounded"
            loading={deletingList}
            color="red"
          />
          <Button
            text="Create new Item"
            loadingText="Creating new Item..."
            onClick={createItem}
            className="inline-flex sm:ml-3 mt-4 sm:mt-0 items-start justify-start px-6 py-3 bg-green-600 hover:bg-green-500 text-white focus:outline-none rounded"
            loading={creatingItem}
          />
          <table className="w-full whitespace-nowrap">
            <thead>
              <tr className="h-16 w-full text-sm leading-none ">
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
              </tr>
            </thead>
            <tbody className="w-full">
              <ReactSortable list={data?.items} setList={setItems}>
                {data?.items.map((e: Item, idx: number) => {
                  return (
                    <tr
                      className="h-20 text-lg leading-none text-gray-800 bg-white hover:bg-gray-100 border-b border-t border-gray-100"
                      key={idx}
                    >
                      <td className="pl-16">
                        <p className="font-medium">{e.title}</p>
                      </td>
                      <td className="pl-16">
                        <p className="font-medium">
                          {e.bought ? (
                            <span className="text-green-500"></span>
                          ) : (
                            <span className="text-red-500">Not Bought</span>
                          )}
                        </p>
                      </td>
                      <td className="pl-16"></td>
                      <td className="pl-1"></td>
                    </tr>
                  );
                })}
              </ReactSortable>
            </tbody>
          </table>
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
