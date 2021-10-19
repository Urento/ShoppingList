import React from "react";
import { useState } from "react";
import { useHistory } from "react-router-dom";
import swal from "sweetalert";
import { queryClient } from "../..";
import { Button } from "../../components/Button";
import { Loading } from "../../components/Loading";
import { Sidebar } from "../../components/Sidebar";
import useAuthCheck from "../../hooks/useAuthCheck";
import { API_URL } from "../../util/constants";

interface CreateResponseData {
  message: string;
  success: "true" | "false";
  owner: string;
  position: number;
  participants: string[];
}

interface CreateResponse {
  message: string;
  data: CreateResponseData;
  code: number;
}

interface Props {
  open: boolean;
}

export const NewShoppinglist: React.FC<Props> = ({ open }) => {
  const [title, setTitle] = useState("");
  const [loading, setLoading] = useState(false);
  const history = useHistory();
  const authStatus = useAuthCheck();

  if (authStatus === "fail") {
    localStorage.removeItem("authenticated");
    history.push("/");
  }

  if (authStatus === "pending") return <Loading withSidebar />;

  const handleTitleChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setTitle(e.target.value);

  const createList = async (e: any) => {
    e.preventDefault();

    if (title.length <= 0)
      return swal({
        icon: "error",
        title: "Error while creating Shopinglist",
        text: "Title has to be longer than 0 Characters",
      });

    setLoading(true);

    const response = await fetch(`${API_URL}/list`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
      body: JSON.stringify({
        title: title,
      }),
    });

    const fJson: CreateResponse = await response.json();
    if (fJson.code !== 200)
      return swal({
        icon: "error",
        title: "Error creating Shoppnglist",
        text: fJson.message,
      });

    queryClient.invalidateQueries("shoppinglists");
    swal({
      icon: "success",
      title: "Successfully created",
      text: "Shoppinglist successfully created",
    });
    setLoading(false);

    //TODO: clear timeout?
    setTimeout(() => {
      swal.close!();
      history.push("/dashboard");
    }, 2000);
  };

  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
      <div className="container mx-auto">
        <div className="flex justify-center px-6 my-12">
          <div className="w-full xl:w-3/4 lg:w-11/12 flex">
            <div className="bg-white p-5 rounded-lg lg:rounded-l-none">
              <form
                className="px-8 pt-6 pb-8 mb-4 bg-white rounded"
                onSubmit={createList}
              >
                <div className="mb-4 md:flex md:justify-between">
                  <div className="mb-4 md:mr-2 md:mb-0">
                    <label
                      className="mb-2 text-sm font-bold text-gray-700"
                      htmlFor="title"
                    >
                      Title
                    </label>
                    <input
                      className="w-full px-3 py-2 text-sm leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline"
                      id="title"
                      type="text"
                      placeholder="Title"
                      onChange={handleTitleChange}
                    />
                  </div>
                </div>
                <div className="mb-6 text-center">
                  <Button
                    text="Create"
                    loadingText="Creating"
                    loading={loading}
                    onClick={createList}
                    type="submit"
                    color="green"
                  />
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default NewShoppinglist;
