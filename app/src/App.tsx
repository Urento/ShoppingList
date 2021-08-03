import React, { useState } from "react";
import "./App.css";
import { AUTH_API_URL } from "./util/constants";
import { Link, Redirect } from "react-router-dom";
import swal from "sweetalert";
import clsx from "clsx";
import jwtDecode from "jwt-decode";
import { useEffect } from "react";
import { Button } from "./components/Button";

interface DataResponse {
  token: string;
}
interface LoginJSONResponse {
  code: string;
  message: "fail" | "ok";
  data: DataResponse;
}

interface JWTPayload {
  email: string;
  secretId: string;
}

//Implement 2FA Screen
const Login: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState({
    email: false,
    password: false,
  });
  const [redirect, setRedirect] = useState(false);
  const [loading, setLoading] = useState(false);
  const loggedIn = Boolean(localStorage.getItem("authenticated"));

  useEffect(() => {
    //redirect for checking if the user is actually authenticated
    if (loggedIn) {
      setRedirect(true);
    }
  }, []);

  const handleEmailChange = (event: React.ChangeEvent<HTMLInputElement>) =>
    setEmail(event.target.value);
  const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) =>
    setPassword(event.target.value);

  const doLogin = async (event: any) => {
    event.preventDefault();
    setLoading(true);

    const f = await fetch(AUTH_API_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      cache: "no-cache",
      credentials: "include",
      body: JSON.stringify({
        email: email,
        password: password,
      }),
    });
    const fJson: LoginJSONResponse = await f.json();

    if (fJson.message === "fail") {
      setError({ email: true, password: true });
    } else if (fJson.message === "ok") {
      const secretId = getSecretIdByJwtToken(fJson.data.token);
      if (secretId == null) {
        swal({
          icon: "error",
          title: "JWT Token couldn't get decoded",
          text: "Error while decoding jwt token! Try again later!",
        });
      }
      console.log(secretId);

      //display modal
      swal({
        icon: "success",
        title: "Successful Login",
        text: "You have been successfully logged in!",
      });

      //update state
      setError({ email: false, password: false });
      setRedirect(true);
      localStorage.setItem("authenticated", "true");
    } else {
      setError({ email: true, password: true });
    }
    //update loading state
    setLoading(false);
  };

  const getSecretIdByJwtToken = (token: string) => {
    const decoded = jwtDecode<JWTPayload>(token);
    return decoded.secretId;
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      {redirect && <Redirect to="/dashboard"></Redirect>}
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
        <form className="mt-8 space-y-6" onSubmit={doLogin}>
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
                  `text-sm sm:text-base relative w-full border rounded placeholder-gray-400 focus:border-indigo-400 focus:outline-none py-2 pr-2 ${
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
            <Button
              showIcon={true}
              loading={loading}
              onClick={doLogin}
              text="Login"
            ></Button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default Login;
