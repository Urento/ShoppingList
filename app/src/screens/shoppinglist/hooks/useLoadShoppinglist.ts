import { useState } from "react";
import { useQuery } from "react-query";
import { useEffect } from "react-router/node_modules/@types/react";
import { ListResponse, ListResponseData } from "../../../types/Shoppinglist";
import { API_URL } from "../../../util/constants";

const useLoadShoppinglist = (id: number) => {
  const [shoppinglist, setShoppinglist] = useState<ListResponseData>({
    id: 0,
    items: [],
    owner: "",
    participants: [],
    title: "",
    created_on: 0,
    deleted_at: 0,
    modified_on: 0,
  });
  const [loadingShoppinglist, setLoadingShoppinglist] = useState<boolean>(true);
  const [error, setError] = useState<string>("");

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
    return fJson.data;
  };

  const loadShoppinglist = async (id: number) => {
    const { isLoading, data, error } = useQuery<ListResponseData, Error>(
      `shoppinglist_${id}`,
      async () => await fetchData(id)
    );
    if (error) setError(error.message);
    if (!isLoading) {
      setShoppinglist(data!);
    }
    setLoadingShoppinglist(false);
  };

  useEffect(() => {
    loadShoppinglist(id);
  }, []);

  return {
    shoppinglist,
    setShoppinglist,
    loadingShoppinglist,
    setLoadingShoppinglist,
    error,
    setError,
  };
};

export default useLoadShoppinglist;
