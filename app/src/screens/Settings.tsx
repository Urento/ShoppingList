import { LockClosedIcon } from "@heroicons/react/solid";
import clsx from "clsx";
import { useEffect } from "react";
import { useState } from "react";
import { useMutation, useQuery } from "react-query";
import { useHistory } from "react-router-dom";
import swal from "sweetalert";
import { queryClient } from "..";
import { Loading } from "../components/Loading";
import { Sidebar } from "../components/Sidebar";
import { ToggleSwitch } from "../components/ToggleSwitch";
import useAuthCheck from "../hooks/useAuthCheck";
import { useFetchUserData } from "../hooks/useFetchUserData";
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

interface TwoFactorAuthenticationData {
  success: "true" | "false";
  message: string;
  status: "true" | "false";
}

interface TwoFactorAuthenticationResponse {
  message: string;
  data: TwoFactorAuthenticationData;
  code: number;
}

interface UpdateUserData {
  success: "true" | "false";
  message: string;
  error: string;
}

interface UpdateUserResponse {
  message: string;
  data: UpdateUserData;
  code: number;
}

export const Settings: React.FC = () => {
  const history = useHistory();
  const authStatus = useAuthCheck();
  const [toggled, setToggled] = useState<boolean | undefined>(false);
  const [email, setEmail] = useState<string | undefined>("");
  const [oldPassword, setOldPassword] = useState<string | undefined>("");
  const [password, setPassword] = useState<string | undefined>("");
  const [confirmPassword, setConfirmPassword] = useState<string | undefined>(
    ""
  );
  const [username, setUsername] = useState<string | undefined>("");
  const [loading, setLoading] = useState<boolean>(false);

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setEmail(e.target.value);
  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setPassword(e.target.value);
  const handleUsernameChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setUsername(e.target.value);
  const handleConfirmPasswordChange = (
    e: React.ChangeEvent<HTMLInputElement>
  ) => setConfirmPassword(e.target.value);
  const handleOldPasswordChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setOldPassword(e.target.value);

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

  const addMutation = useMutation(
    () =>
      fetch(`${API_URL}/auth/user`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        credentials: "include",
      }),
    {
      onSuccess: () => queryClient.invalidateQueries("user"),
    }
  );

  useEffect(() => {
    if (!isLoading) {
      setEmail(data?.data.e_mail);
      setUsername(data?.data.username);
      setToggled(data?.data.two_factor_authentication);
    }
  }, [isLoading]);

  if (authStatus === "fail") {
    localStorage.removeItem("authenticated");
    history.push("/");
  }

  if (authStatus === "pending") {
    return <Loading />;
  }

  //fetch user

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    swal({
      icon: "error",
      title: "Error while fetching settings",
      text: "Error while getting the user information! Try again later",
    });
    history.push("/dashboard");
  }

  const updateTwoFactorAuthentication = async () => {
    const response = await fetch(`${API_URL}/twofactorauthentication`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
      body: JSON.stringify({
        status: !toggled,
        user_token: "none", //TODO
      }),
    });
    const fJson: TwoFactorAuthenticationResponse = await response.json();
    if (fJson.code != 200) {
      swal({
        icon: "error",
        title: `Error ${
          toggled ? "disabling" : "enabling"
        } Two Facotor Authentication`,
        text: `An error occurred while ${
          toggled ? "disabling" : "enabling"
        } Two Factor Authentication!`,
      });
      return;
    }
    swal({
      icon: "success",
      title: `Two Factor Authentication ${
        fJson.data.status === "true" ? "enabled" : "disabled"
      }!`,
      text: `You successfully ${
        fJson.data.status === "true" ? "enabled" : "disabled"
      } Two Factor Authentication`,
    });
    setToggled(fJson.data.status === "true");
  };

  const updateUser = async (e: any) => {
    e.preventDefault();
    setLoading(true);

    if (password != confirmPassword) {
      swal({
        icon: "error",
        title: "Password and Confirm Password have to be the same",
        text: "Please check your password and confirm password again!",
      });
      setLoading(false);
      return;
    }

    const response = await fetch(`${API_URL}/auth/update`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
      body: JSON.stringify({
        email: email,
        username: username,
        password: password,
        with_password: password != "" ? true : false,
        old_password: oldPassword,
      }),
    });
    const fJson: UpdateUserResponse = await response.json();
    if (fJson.code != 200) {
      swal({
        icon: "error",
        title: "Error while updating profile!",
        text: fJson.data != null ? fJson.data.message : fJson.message,
      });
      return;
    }
    swal({
      icon: "success",
      title: "Successfully updated!",
      text: "Successfully updated your account!",
    });

    if (error) {
      history.push("/dashboard");
      return;
    }

    addMutation.mutate();

    if (!isLoading && error) {
      setEmail(data?.data.e_mail);
      setUsername(data?.data.username);
      setToggled(data?.data.two_factor_authentication);
      //maybe clear password fields?
    }
    setLoading(false);
  };

  //TODO: FIX DESIGN: Looks horrible but it does the job
  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
      <div className="container mx-auto">
        <div className="flex justify-center px-6 my-12">
          <div className="w-full xl:w-3/4 lg:w-11/12 flex">
            <div className="w-full lg:w-7/12 bg-white p-5 rounded-lg lg:rounded-l-none">
              <form
                className="px-8 pt-6 pb-8 mb-4 bg-white rounded"
                onSubmit={updateUser}
              >
                <div className="mb-4 md:flex md:justify-between">
                  <div className="mb-4 md:mr-2 md:mb-0">
                    <label
                      className="mb-2 text-sm font-bold text-gray-700"
                      htmlFor="email"
                    >
                      Email
                    </label>
                    <input
                      className="w-full px-3 py-2 text-sm leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline"
                      id="email"
                      type="email"
                      placeholder="Email address"
                      value={email}
                      onChange={handleEmailChange}
                    />
                  </div>
                  <div className="md:ml-2">
                    <label
                      className="mb-2 text-sm font-bold text-gray-700"
                      htmlFor="username"
                    >
                      Username
                    </label>
                    <input
                      className="w-full px-3 py-2 text-sm leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline"
                      id="username"
                      type="text"
                      placeholder="Username"
                      value={username}
                      onChange={handleUsernameChange}
                    />
                  </div>
                </div>
                <div className="mb-4">
                  <ToggleSwitch
                    id="twoFactorAuthentication"
                    onClick={updateTwoFactorAuthentication}
                    toggled={toggled!}
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
                      className="w-full px-3 py-2 mb-3 text-sm leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline"
                      id="password"
                      type="password"
                      onChange={handleOldPasswordChange}
                    />
                  </div>
                  <div className="md:ml-2">
                    <label
                      className="block mb-2 text-sm font-bold text-gray-700"
                      htmlFor="password"
                    >
                      New Password
                    </label>
                    <input
                      className="w-full px-3 py-2 mb-3 text-sm leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline"
                      id="password"
                      type="password"
                      onChange={handlePasswordChange}
                    />
                  </div>
                  <div className="md:ml-2">
                    <label
                      className="block mb-2 text-sm font-bold text-gray-700"
                      htmlFor="confirmPassword"
                    >
                      Confirm Password
                    </label>
                    <input
                      className="w-full px-3 py-2 mb-3 text-sm leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline"
                      id="confirmPassword"
                      type="password"
                      onChange={handleConfirmPasswordChange}
                    />
                  </div>
                </div>
                <div className="mb-6 text-center">
                  <button
                    className="w-full px-4 py-2 font-bold text-white bg-green-500 rounded-full hover:bg-green-700 focus:outline-none focus:shadow-outline"
                    type="submit"
                    onSubmit={updateUser}
                  >
                    {loading ? (
                      <svg
                        className="loading-svg justify-center flex"
                        viewBox="25 25 50 50"
                      >
                        <circle
                          className="loading-circle"
                          cx="50"
                          cy="50"
                          r="20"
                        ></circle>
                      </svg>
                    ) : (
                      "Update"
                    )}
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
