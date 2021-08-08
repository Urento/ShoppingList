import React, { useState } from "react";
import "./App.css";
import { AUTH_API_URL, TOTP_API_URL } from "./util/constants";
import { Link, Redirect, useHistory, useLocation } from "react-router-dom";
import swal from "sweetalert";
import clsx from "clsx";
import jwtDecode from "jwt-decode";
import { useEffect } from "react";
import { Button } from "./components/Button";
import { isLoggedIn } from "./storage/UserStorage";

interface TOTPDataResponse {
  token: string;
  success: "true" | "false";
  otp: boolean;
  error: string;
  verified: "true" | "false";
}
interface TOTPJSONResponse {
  code: string;
  message: "fail" | "ok";
  data: TOTPDataResponse;
}

interface JWTPayload {
  email: string;
  secretId: string;
}

interface Props {
  email: string;
}

export const TwoFactorAuthentication: React.FC = () => {
  const [otp, setOtp] = useState("");
  const [error, setError] = useState(false);
  const [loading, setLoading] = useState(false);
  const loggedIn = isLoggedIn();
  const history = useHistory();
  const location = useLocation<Props>();

  useEffect(() => {
    //redirect for checking if the user is actually already authenticated
    if (loggedIn) {
      history.push("/dashboard");
    }

    if (
      location.state.email === null ||
      location.state.email === "" ||
      location.state.email === undefined
    )
      history.push("/");
  }, []);

  const handleOTPChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setOtp(e.target.value.replace(" ", ""));

  const verifyOTP = async (event: any) => {
    event.preventDefault();
    setLoading(true);

    const f = await fetch(TOTP_API_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      cache: "no-cache",
      credentials: "include",
      body: JSON.stringify({
        otp: otp,
        login_after: true,
        email: location.state.email,
        enable_after: false,
      }),
    });
    const fJson: TOTPJSONResponse = await f.json();
    console.log(fJson.data);

    if (fJson.message === "fail") {
      setError(true);
    } else if (
      fJson.message === "ok" &&
      fJson.data.verified === "true" &&
      fJson.data.success === "true"
    ) {
      const secretId = getSecretIdByJwtToken(fJson.data.token);
      if (secretId == null) {
        swal({
          icon: "error",
          title: "JWT Token couldn't get decoded",
          text: "Error while decoding jwt token! Try again later!",
        });
      }

      //display modal
      swal({
        icon: "success",
        title: "Successful Login",
        text: "You have been successfully logged in!",
      });

      //update state
      setError(false);
      history.push("/dashboard");
      localStorage.setItem("authenticated", "true");
    } else {
      setError(true);
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
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Two Factor Authentication
          </h2>
        </div>
        <form className="mt-8 space-y-6" onSubmit={verifyOTP}>
          <input type="hidden" name="remember" defaultValue="true" />
          <div className="rounded-md shadow-sm -space-y-px">
            <div>
              <label htmlFor="email-address" className="sr-only">
                OTP
              </label>
              <input
                id="email-address"
                name="email"
                type="number"
                autoComplete="email"
                required
                className={clsx(
                  `text-sm sm:text-base relative w-full border rounded placeholder-gray-400 focus:border-indigo-400 focus:outline-none py-2 pr-2 ${
                    error ? "border-red-500" : ""
                  }`
                )}
                onChange={handleOTPChange}
              />
              {error ? (
                <span className="flex items-center font-medium tracking-wide text-red-500 text-xs mt-1 ml-1">
                  Invalid OTP!
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
              onClick={verifyOTP}
              text="Verify"
              type="submit"
              danger={false}
            ></Button>
          </div>
        </form>
      </div>
    </div>
  );
};
