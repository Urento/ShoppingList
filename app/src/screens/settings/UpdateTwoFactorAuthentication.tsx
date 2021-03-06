import React, { useState } from "react";
import { useEffect } from "react";
import { useHistory, useLocation } from "react-router-dom";
import swal from "sweetalert";
import { queryClient } from "../..";
import { Button } from "../../components/Button";
import { Loading } from "../../components/Loading";
import { Sidebar } from "../../components/Sidebar";
import useAuthCheck from "../../hooks/useAuthCheck";
import { getUser } from "../../storage/UserStorage";
import { VerifyResponse } from "../../types/TwoFactorAuthentication";
import { API_URL, TOTP_API_URL } from "../../util/constants";
import { useLoadQRCode } from "./hooks/useLoadQRCode";

interface Props {
  status: boolean;
}

const UpdateTwoFactorAuthentication: React.FC = () => {
  const location = useLocation<Props>();
  const qrCode = useLoadQRCode(location.state.status);
  const history = useHistory();
  const authStatus = useAuthCheck();
  const [otp, setOTP] = useState("");
  const [status, setStatus] = useState(true); //enable or disable totp
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (location.state.status === null || location.state.status === undefined)
      history.push("/");
    setStatus(location.state.status);
    console.log(location.state.status);
  }, [location]);

  if (authStatus === "fail") {
    localStorage.removeItem("authenticated");
    history.push("/");
  }

  if (authStatus === "pending") return <Loading withSidebar />;

  const enableTOTP = async (e: any) => {
    e.preventDefault();
    if (!status) return; //maybe display error message
    setLoading(true);
    const user = await getUser();
    const email = user.email;

    const response = await fetch(TOTP_API_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
      body: JSON.stringify({
        email: email,
        otp: otp,
        login_after: false,
        enable_after: true,
      }),
    });
    const fJson: VerifyResponse = await response.json();
    if (
      fJson.code !== 200 ||
      fJson.data.success !== "true" ||
      fJson.data.verified !== "true"
    ) {
      setLoading(false);
      return swal({
        icon: "error",
        title: "Error while enabling TOTP",
        text: "Please contact an administrator!",
      });
    }

    queryClient.invalidateQueries("user");
    setTimeout(() => {
      history.push("/");
    }, 5000);
    swal({
      icon: "success",
      title: "Successfully activated!",
      text: "You will be redirected in 5 seconds!",
    });
    setLoading(false);
    //TODO: Display Backup Codes
  };

  const verifyTOTP = async (e: any) => {
    e.preventDefault();
    if (status) return; //maybe display error message
    setLoading(true);

    const response = await fetch(`${API_URL}/twofactorauthentication`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
      body: JSON.stringify({
        otp: otp,
        status: false,
      }),
    });
    const fJson: VerifyResponse = await response.json();
    if (
      fJson.code !== 200 ||
      fJson.data.success !== "true" ||
      fJson.data.verified !== "true"
    ) {
      setLoading(false);
      return swal({
        icon: "error",
        title: "Error disabling Two Factor Authentication",
      });
    }

    queryClient.invalidateQueries("user");
    setTimeout(() => {
      history.push("/");
    }, 5000);
    swal({
      icon: "success",
      title: "Successfully disbaled Two Factor Authentication",
      text: "You will be redirected in 5 seconds!",
    });
    setLoading(false);
  };

  const handleOTPChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setOTP(e.target.value);

  //TODO: display button to disable it again if you enabled it by accident

  if (qrCode === "Error while loading QRCode" && location.state.status) {
    return <div>Unable to load QRCode! Try again later!</div>;
  }

  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
      <div className="container mx-auto">
        <div className="flex justify-center px-6 my-12">
          <div className="w-full xl:w-3/4 lg:w-11/12 flex">
            <h1 className="text-2xl">
              {status ? "Enable" : "Disable"} Two Factor Authentication
            </h1>
            {status && (
              <img alt="QR Code" src={`data:image/png;base64,${qrCode}`} />
            )}
            <form
              className="px-8 pt-6 pb-8 mb-4 bg-white rounded"
              onSubmit={location.state.status ? enableTOTP : verifyTOTP}
            >
              <div className="mb-4 md:flex md:justify-between">
                <div className="mb-4 md:mr-2 md:mb-0">
                  <label
                    className="mb-2 text-sm font-bold text-gray-700"
                    htmlFor="otp"
                  >
                    OTP
                  </label>
                  <input
                    className="w-full px-3 py-2 text-sm leading-tight text-gray-700 border rounded shadow appearance-none focus:outline-none focus:shadow-outline"
                    id="otp"
                    type="number"
                    placeholder="OTP"
                    onChange={handleOTPChange}
                  />
                </div>
              </div>
              {status && (
                <Button
                  loading={loading}
                  onClick={enableTOTP}
                  showIcon={true}
                  text="Enable"
                  type="submit"
                  danger={false}
                />
              )}
              <br />
              {!status && (
                <Button
                  loading={loading}
                  onClick={enableTOTP}
                  showIcon={true}
                  text="Disable"
                  type="submit"
                  danger={true}
                />
              )}
            </form>
          </div>
        </div>
      </div>
    </div>
  );
};

export default UpdateTwoFactorAuthentication;
