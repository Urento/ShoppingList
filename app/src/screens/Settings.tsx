import clsx from "clsx";
import { useState } from "react";
import { useQuery } from "react-query";
import { useHistory } from "react-router-dom";
import swal from "sweetalert";
import { Button } from "../components/Button";
import { Loading } from "../components/Loading";
import { Sidebar } from "../components/Sidebar";
import { ToggleSwitch } from "../components/ToggleSwitch";
import useAuthCheck from "../hooks/useAuthCheck";
import { API_URL } from "../util/constants";

interface UserInfo {
  id: number;
  e_mail: string;
  email_verified: boolean;
  username: string;
  rank: string;
  two_factor_authentication: boolean;
}

interface UserInfoResponse {
  data: UserInfo;
  message: string;
  code: string;
}

export const Settings: React.FC = () => {
  const history = useHistory();
  const authStatus = useAuthCheck();
  const [toggled, setToggled] = useState(false);

  const { isLoading, error, data } = useQuery<UserInfoResponse, Error>(
    "user",
    () =>
      fetch(`${API_URL}/auth/user`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        credentials: "include",
      }).then((res) => res.json()),
    { refetchOnWindowFocus: false }
  );

  if (authStatus === "fail") {
    localStorage.removeItem("authenticated");
    history.push("/");
  }

  if (authStatus === "pending") {
    return <Loading />;
  }

  //fetch user

  if (isLoading) {
    return <>Loading...</>;
  }

  if (error) {
    swal({
      icon: "error",
      title: "Error while fetching settings",
      text: "Error while getting user information! Try again later",
    });
    history.push("/dashboard");
  }

  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
      <div className="container mx-auto">
        <div className="flex justify-center px-6 my-12">
          <div className="w-full xl:w-3/4 lg:w-11/12 flex">
            <div className="w-full lg:w-7/12 bg-white p-5 rounded-lg lg:rounded-l-none">
              <form className="px-8 pt-6 pb-8 mb-4 bg-white rounded">
                <div className="mb-4 md:flex md:justify-between">
                  <div className="mb-4 md:mr-2 md:mb-0">
                    <label
                      className="block mb-2 text-sm font-bold text-gray-700"
                      htmlFor="email"
                    >
                      Email
                    </label>
                    <input
                      className="w-full px-9 py-2 text-sm leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline"
                      id="email"
                      type="email"
                      placeholder="Email address"
                      value={data?.data.e_mail}
                    />
                  </div>
                  <div className="md:ml-2">
                    <label
                      className="block mb-2 text-sm font-bold text-gray-700"
                      htmlFor="username"
                    >
                      Username
                    </label>
                    <input
                      className="w-full px-3 py-2 text-sm leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline"
                      id="username"
                      type="text"
                      placeholder="Username"
                      value={data?.data.username}
                    />
                  </div>
                </div>
                <div className="mb-4">
                  <ToggleSwitch
                    id="twoFactorAuthentication"
                    onClick={() => setToggled(!toggled)}
                    toggled={data?.data.two_factor_authentication!}
                    title="Two Factor Authentication"
                  />
                </div>
                <div className="mb-4 md:flex md:justify-between">
                  <div className="mb-4 md:mr-2 md:mb-0">
                    <label
                      className="block mb-2 text-sm font-bold text-gray-700"
                      htmlFor="password"
                    >
                      Old Password
                    </label>
                    <input
                      //border-red-500 on error
                      className="w-full px-3 py-2 mb-3 text-sm leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline"
                      id="password"
                      type="password"
                      placeholder="******************"
                    />
                    {/*<p className="text-xs italic text-red-500">
                      Please choose a password.
  </p>*/}
                  </div>
                  <div className="md:ml-2">
                    <label
                      className="block mb-2 text-sm font-bold text-gray-700"
                      htmlFor="c_password"
                    >
                      New Password
                    </label>
                    <input
                      className="w-full px-3 py-2 mb-3 text-sm leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline"
                      id="c_password"
                      type="password"
                      placeholder="******************"
                    />
                  </div>
                  <div className="md:ml-2">
                    <label
                      className="block mb-2 text-sm font-bold text-gray-700"
                      htmlFor="c_password"
                    >
                      Confirm Password
                    </label>
                    <input
                      className="w-full px-3 py-2 mb-3 text-sm leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline"
                      id="c_password"
                      type="password"
                      placeholder="******************"
                    />
                  </div>
                </div>
                <div className="mb-6 text-center">
                  <button
                    className="w-full px-4 py-2 font-bold text-white bg-green-500 rounded-full hover:bg-blue-700 focus:outline-none focus:shadow-outline"
                    type="button"
                  >
                    Update
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
