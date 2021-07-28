import React, { useState } from "react";
import "./App.css";
import { LockClosedIcon } from "@heroicons/react/solid";
import { AUTH_API_URL } from "./util/constants";
import { Link } from "react-router-dom";
import swal from "sweetalert";
import clsx from "clsx";

interface DataResponse {
  token: string;
}
interface LoginJSONResponse {
  code: string;
  message: "fail" | "ok";
  data: DataResponse;
}

const checkLogin = async () => {
  const token = localStorage.getItem("token");
  if (token != null && token!.length > 0) {
    //TODO: send request to backend to check auth token
  }
};

const Login: React.FC = () => {
  checkLogin();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState({
    email: false,
    password: false,
  });

  const handleEmailChange = (event: any) => setEmail(event.target.value);
  const handlePasswordChange = (event: any) => setPassword(event.target.value);

  const login = async (event: any) => {
    event.preventDefault();
    const f = await fetch(AUTH_API_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email: email,
        password: password,
      }),
    });
    const fJson: LoginJSONResponse = await f.json();
    if (fJson.message === "fail") {
      setError({ email: true, password: true });
    } else if (fJson.message === "ok") {
      swal({
        icon: "success",
        title: "Successful Login",
        text: "You have been successfully logged in!",
      });
    } else {
      setError({ email: true, password: true });
    }
    //TODO: do redux stuff
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Login
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            Or{" "}
            <Link
              to="/register"
              className="font-medium text-indigo-600 hover:text-indigo-500"
            >
              Create a new Account!
            </Link>
          </p>
        </div>
        <form className="mt-8 space-y-6" onSubmit={login}>
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
                  `text-sm sm:text-base relative w-full border rounded placeholder-gray-400 focus:border-indigo-400 focus:outline-none py-2 pr-2 pl-12 ${
                    error.email ? "border-red-500" : ""
                  }`
                )}
                placeholder="Email address"
                onChange={handleEmailChange}
              />
              {error.email ? (
                <span className="flex items-center font-medium tracking-wide text-red-500 text-xs mt-1 ml-1">
                  Invalid email or password!
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
                  `text-sm sm:text-base relative w-full border rounded placeholder-gray-400 focus:border-indigo-400 focus:outline-none py-2 pr-2 pl-12 ${
                    error.password ? "border-red-500" : ""
                  }`
                )}
                placeholder="Password"
                onChange={handlePasswordChange}
              />
            </div>
            {error.password ? (
              <span className="flex items-center font-medium tracking-wide text-red-500 text-xs mt-1 ml-1">
                Invalid email or password!
              </span>
            ) : (
              ""
            )}
          </div>

          <div className="flex items-center justify-between">
            <div className="text-sm">
              <a
                href="#"
                className="font-medium text-indigo-600 hover:text-indigo-500"
              >
                Forgot your password?
              </a>
            </div>
          </div>

          <div>
            <button
              type="submit"
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              onSubmit={login}
            >
              <span className="absolute left-0 inset-y-0 flex items-center pl-3">
                <LockClosedIcon
                  className="h-5 w-5 text-indigo-500 group-hover:text-indigo-400"
                  aria-hidden="true"
                />
              </span>
              Login
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default Login;
