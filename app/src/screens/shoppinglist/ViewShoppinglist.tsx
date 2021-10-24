import React from "react";
import { useQuery } from "react-query";
import { useHistory, useParams } from "react-router";
import { useState, useEffect } from "react";
import { Button } from "../../components/Button";
import { Sidebar } from "../../components/Sidebar";
import { API_URL } from "../../util/constants";
import { Loading } from "../../components/Loading";
import { ReactSortable, Sortable } from "react-sortablejs";
import swal from "sweetalert";
import {
  Item,
  ItemType,
  ListResponse,
  ListResponseData,
} from "../../types/Shoppinglist";
import { DeleteResponse } from "../../components/ShoppinglistCard";
import { configureStore } from "@reduxjs/toolkit";

interface Params {
  id: string;
}

interface DeletingItemState {
  id: number;
  loading: boolean;
}

interface MarkingItemAsBoughtState {
  id: number;
  loading: boolean;
}

//TODO: Implement Auth Hook
const Viewdata: React.FC = ({}) => {
  const { id } = useParams<Params>();
  const history = useHistory();
  const [creatingItem, setCreatingItem] = useState(false);
  const [deletingList, setDeletingList] = useState(false);
  const [deletingItem, setDeletingItem] = useState<DeletingItemState>({
    id: 0,
    loading: false,
  });
  const [items, setItems] = useState<Item[]>([]);

  const { isLoading, data, error, refetch, isFetching } = useQuery<
    ListResponse,
    Error
  >(`data_${id}`, async () => await fetchData(parseInt(id)), {
    refetchOnWindowFocus: false,
  });

  useEffect(() => {
    //TODO: Fix sorting
    const sortedArray: Item[] | undefined = data?.data.items.sort(
      (item1, item2) => {
        if (item1.position > item2.position) return 1;
        if (item1.position < item2.position) return -1;
        return 0;
      }
    );

    if (!isLoading) setItems(sortedArray!);
    if (!isFetching) setItems(sortedArray!);
  }, [isLoading, isFetching]);

  if (error) return <Loading withSidebar />;
  if (isLoading) return <Loading withSidebar />;
  if (!data) return <Loading withSidebar />;
  if (isFetching) return <Loading withSidebar />;

  const createItem = async () => {
    setCreatingItem(true);
    let pos: number;
    let position: number;

    try {
      pos = data?.data.items[data?.data.items.length - 1].position;
      position = pos + 1;
    } catch {
      position = 1;
    }

    await swal({
      text: "Create a new Item", //@ts-ignore
      content: "input",
      buttons: ["Cancel", "Create"],
      closeOnEsc: false,
      closeOnClickOutside: false,
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
      setCreatingItem(false);
    });
  };

  const deleteList = async () => {
    setDeletingList(true);
    const alert = await swal({
      icon: "warning",
      title: "Are you sure you want to delete the Shoppinglist?",
      buttons: ["No", "Yes"],
    });

    if (alert) {
      const response = await fetch(`${API_URL}/list/${id}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        credentials: "include",
      });
      const fJson: DeleteResponse = await response.json();

      if (fJson.code !== 200)
        return swal({
          icon: "error",
          title: "Error while deleting",
          text: "Error while deleting Shoppinglist. Try again later!",
        });

      swal({
        icon: "success",
        title: "Successfully deleted",
        text: "Successfully deleted the Shoppinglist",
      });

      setTimeout(() => {
        history.push("/");
      }, 1000);
    }

    setDeletingList(false);
  };

  const deleteItem = async (itemId: number) => {
    setDeletingItem({ id: itemId, loading: true });

    await fetch(`${API_URL}/item`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify({ id: itemId, parent_list_id: parseInt(id) }),
      credentials: "include",
    });
    refetch();
    setDeletingItem({ id: 0, loading: false });
  };

  const updateItems = async (evt: Sortable.SortableEvent) => {
    console.log(evt.oldIndex! + 1);
    console.log(evt.newIndex! + 1);
    const item1 = items[evt.oldIndex!];
    const item2 = items[evt.newIndex!];
    const i: ItemType[] = [
      {
        title: item1.title,
        bought: item1.bought,
        id: item1.id,
        itemId: item1.itemId,
        parentListId: item1.parentListId,
        position: evt.newIndex! + 1,
      },
      {
        title: item2.title,
        bought: item2.bought,
        id: item2.id,
        itemId: item2.itemId,
        parentListId: item2.parentListId,
        position: evt.oldIndex! + 1,
      },
    ];
    await fetch(`${API_URL}/items`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify({ parent_list_id: parseInt(id), items: i }),
      credentials: "include",
    });
    refetch();
  };

  const updateItem = async (item: ItemType) => {
    await fetch(`${API_URL}/item/${item.itemId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify({
        id: item.id,
        parentListId: item.parentListId,
        title: item.title,
        position: item.position,
        bought: item.bought,
      }),
      credentials: "include",
    });

    refetch();
  };

  const leaveShoppinglist = async () => {
    const response = await fetch(`${API_URL}/participant/list/leave`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify({ id: parseInt(id) }),
      credentials: "include",
    });

    if (response) setTimeout(() => history.push("/"), 500);
  };

  const updateTitle = async (item: ItemType) => {
    await swal({
      text: "Update Title", //@ts-ignore
      content: "input",
      buttons: ["Cancel", "Update"],
      closeOnEsc: false,
      closeOnClickOutside: false,
    }).then(async (title: string) => {
      if (title === "" || !title) return;
      const response = await fetch(`${API_URL}/item/${item.itemId}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: JSON.stringify({
          parentListId: item.parentListId,
          title: title,
          bought: item.bought,
          position: item.position,
        }),
        credentials: "include",
      });

      if (response) refetch();
    });
  };

  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
      <div className="container mx-auto py-10 md:w-4/5 w-11/12">
        <div className="bg-white px-4 md:px-10 pt-5 md:pt-7 pb-5 overflow-y-auto">
          <h1 className="text-2xl">{data?.data.title}</h1>
          <h1 className="text-2xl">
            {items.length} {items.length === 1 ? "Item" : "Items"}
          </h1>
          <h1 className="text-2xl">
            {data?.data.participants.length === 0
              ? 1
              : data?.data.participants.length}{" "}
            {data?.data.participants.length === 1 ||
            data?.data.participants.length === 0
              ? "Participant"
              : "Participants"}
          </h1>
          <br />
          {!data?.is_participant && (
            <Button
              text="Delete Shoppinglist"
              loadingText="Deleting Shoppinglist..."
              onClick={deleteList}
              className="inline-flex sm:ml-3 mt-4 sm:mt-0 items-start justify-start px-6 py-3 bg-red-600 hover:bg-red-500 text-white focus:outline-none rounded"
              loading={deletingList}
              color="red"
            />
          )}
          <Button
            text="Create new Item"
            loadingText="Creating new Item..."
            onClick={createItem}
            className="inline-flex sm:ml-3 mt-4 sm:mt-0 items-start justify-start px-6 py-3 bg-green-600 hover:bg-green-500 text-white focus:outline-none rounded"
            loading={creatingItem}
          />
          {!data?.is_participant && (
            <Button
              text="Participants"
              onClick={() => history.push(`/list/participants/${id}`)}
              className="inline-flex sm:ml-3 mt-4 sm:mt-0 items-start justify-start px-6 py-3 bg-green-600 hover:bg-green-500 text-white focus:outline-none rounded"
            />
          )}
          {data?.is_participant && (
            <Button
              text="Leave"
              onClick={() => leaveShoppinglist()}
              className="inline-flex sm:ml-3 mt-4 sm:mt-0 items-start justify-start px-6 py-3 bg-red-600 hover:bg-red-500 text-white focus:outline-none rounded"
            />
          )}
          <table className="w-full whitespace-nowrap">
            <thead>
              <tr className="h-16 w-full text-sm leading-none ">
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
              </tr>
            </thead>
            <tbody className="w-full">
              <ReactSortable
                list={items}
                setList={setItems}
                onUpdate={updateItems}
              >
                {items.map((e: Item, idx: number) => {
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
                            <span className="text-green-500">Bought</span>
                          ) : (
                            <span className="text-red-500">Not Bought</span>
                          )}
                        </p>
                      </td>
                      <td className="pl-16">
                        <Button
                          text="Delete"
                          loadingText="Deleting..."
                          loading={
                            deletingItem.id === e.itemId && deletingItem.loading
                          }
                          color="red"
                          onClick={() => deleteItem(e.itemId)}
                        />
                      </td>
                      <td className="pl-16">
                        <Button
                          text="Bought"
                          color="green"
                          onClick={() =>
                            updateItem({
                              id: e.id,
                              bought: !e.bought,
                              itemId: e.itemId,
                              parentListId: e.parentListId,
                              position: e.position,
                              title: e.title,
                            })
                          }
                        />
                      </td>
                      <td className="pl-16">
                        <Button
                          text="Edit"
                          color="green"
                          onClick={() =>
                            updateTitle({
                              id: e.id,
                              itemId: e.itemId,
                              bought: e.bought,
                              parentListId: e.parentListId,
                              position: e.position,
                              title: e.title,
                            })
                          }
                        />
                      </td>
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

const fetchData = async (id: number) => {
  const response = await fetch(`${API_URL}/list/${id}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    credentials: "include",
  });
  const fJson: ListResponse = await response.json();
  return fJson;
};

export default Viewdata;
