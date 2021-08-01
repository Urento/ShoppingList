import React, { useState } from "react";
import "./App.css";
import { AUTH_REGISTER_API_URL } from "./util/constants";
import swal from "sweetalert";
import clsx from "clsx";
import { Redirect } from "react-router-dom";
import { Button } from "./components/Button";

interface DataResponse {
  created: "true" | "false";
  email: string;
  username: string;
  error: string;
  success: "true" | "false";
}

interface LoginJSONResponse {
  code: string;
  message: string;
  data: DataResponse;
}

//TODO: Add more custom error cases

const Register: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [username, setUsername] = useState("");
  const [error, setError] = useState({
    email: false,
    username: false,
    password: false,
  });
  const [redirect, setRedirect] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleEmailChange = (event: React.ChangeEvent<HTMLInputElement>) =>
    setEmail(event.target.value);
  const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) =>
    setPassword(event.target.value);
  const handleUsernameChange = (event: React.ChangeEvent<HTMLInputElement>) =>
    setUsername(event.target.value);

  const register = async (event: any) => {
    event.preventDefault();
    setLoading(true);

    const f = await fetch(AUTH_REGISTER_API_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      cache: "no-cache",
      credentials: "include",
      body: JSON.stringify({
        email: email,
        username: username,
        password: password,
      }),
    });
    const fJson: LoginJSONResponse = await f.json();
    if (fJson.data.error === "email is already being used") {
      swal({
        icon: "error",
        title: "Email is already being used!",
        text: "Please try another email!",
      });
    } else if (
      fJson.data.error != null &&
      fJson.data.error.includes("32 characters")
    ) {
      swal({
        icon: "error",
        title: "Username has to be shorter!",
        text: "Username can't be longer than 32 characters!",
      });
    } else if (fJson.message === "fail") {
      setError({ email: true, password: true, username: true });
    } else if (fJson.message === "ok") {
      //display modal
      swal({
        icon: "success",
        title: "Successfully created your account!",
        text: "You have successfully created your account!",
      });

      //update state
      setError({ email: false, password: false, username: false });
      setRedirect(true);
    } else {
      setError({ email: true, password: true, username: true });
    }
    //update loading state
    setLoading(false);
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        {redirect && <Redirect to="/"></Redirect>}
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Register
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            Or{" "}
            <a
              href="/"
              className="font-medium text-indigo-600 hover:text-indigo-500"
            >
              Login into an already existing account!
            </a>
          </p>
        </div>
        <form className="mt-8 space-y-6" onSubmit={register}>
          <input type="hidden" name="remember" defaultValue="true" />
          <div className="rounded-md shadow-sm -space-y-px">
            <div>
              <label htmlFor="email-address" className="sr-only">
                Email address
              </label>
              <input
                id="email-address"
                name="email"
                type="email"
                autoComplete="email"
                required
                className={clsx(
                  `text-sm sm:text-base relative w-full border rounded placeholder-gray-400 focus:border-indigo-400 focus:outline-none py-2 pr-2 ${
                    error.email ? "border-red-500" : ""
                  }`
                )}
                placeholder="Email address"
                onChange={handleEmailChange}
              />
              {error.email ? (
                <span className="flex items-center font-medium tracking-wide text-red-500 text-xs mt-1 ml-1">
                  Invalid email, username or password!
                </span>
              ) : (
                ""
              )}
            </div>

            <div style={{ paddingTop: "2%" }}>
              <label htmlFor="username" className="sr-only">
                Username
              </label>
              <input
                id="username"
                name="username"
                type="text"
                autoComplete="username"
                className={clsx(
                  `text-sm sm:text-base relative w-full border rounded placeholder-gray-400 focus:border-indigo-400 focus:outline-none py-2 pr-2 ${
                    error.username ? "border-red-500" : ""
                  }`
                )}
                placeholder="Username"
                onChange={handleUsernameChange}
              />
              {error.email ? (
                <span className="flex items-center font-medium tracking-wide text-red-500 text-xs mt-1 ml-1">
                  Invalid email, username or password!
                </span>
              ) : (
                ""
              )}
            </div>
            <div style={{ paddingTop: "2%" }}>
              <label htmlFor="password" className="sr-only">
                Password
              </label>
              <input
                id="password"
                name="password"
                type="password"
                autoComplete="current-password"
                required
                className={clsx(
                  `text-sm sm:text-base relative w-full border rounded placeholder-gray-400 focus:border-indigo-400 focus:outline-none py-2 pr-2 ${
                    error.email ? "border-red-500" : ""
                  }`
                )}
                placeholder="Password"
                onChange={handlePasswordChange}
              />
              {error.password ? (
                <span className="flex items-center font-medium tracking-wide text-red-500 text-xs mt-1 ml-1">
                  Invalid email, username or password!
                </span>
              ) : (
                ""
              )}
            </div>
          </div>
          <div>
            <Button
              showIcon={true}
              loading={loading}
              text="Create Account"
              onClick={register}
            />
          </div>
        </form>
      </div>
    </div>
  );
};

export default Register;
