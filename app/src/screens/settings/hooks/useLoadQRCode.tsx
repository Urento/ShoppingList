import { useEffect, useState } from "react";
import { API_URL } from "../../../util/constants";

interface EnableTOTPResponseData {
  success: "true" | "false";
  img: string;
}

interface EnableTOTPResponse {
  code: number;
  data: EnableTOTPResponseData;
  message: string;
}

export const useLoadQRCode = (status = false) => {
  const [qrCode, setQRCode] = useState("");

  useEffect(() => {
    const loadData = async () => {
      const response = await fetch(`${API_URL}/twofactorauthentication`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ otp: "", status: true }),
      });
      const fJson: EnableTOTPResponse = await response.json();
      if (fJson.code !== 200) return setQRCode("Error while loading QRCode");
      if (fJson.data.success !== "true")
        return setQRCode("Error while loading QRCode");
      setQRCode(fJson.data.img);
    };

    if (status) loadData();
  }, [status]);

  return qrCode;
};
