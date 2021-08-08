import { useState, useEffect } from "react";
import { API_URL } from "../util/constants";

interface DataResponse {
  success: "true" | "false";
  error: string;
}

interface AuthCheckResponse {
  code: string;
  message: string;
  data: DataResponse;
}

const useAuthCheck = () => {
  const [status, setStatus] = useState<"success" | "fail" | "pending">(
    "pending"
  );

  useEffect(() => {
    const checkAuth = async () => {
      const response = await fetch(`${API_URL}/auth/check`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        cache: "no-cache",
        credentials: "include",
      });
      const fJson: AuthCheckResponse = await response.json();

      if (
        fJson.data.success !== "true" ||
        fJson.message === "fail" ||
        fJson.message === "not authorized to access this route"
      ) {
        return setStatus("fail");
      } else {
        return setStatus("success");
      }
    };

    checkAuth();
  }, []);

  return status;
};

export default useAuthCheck;
