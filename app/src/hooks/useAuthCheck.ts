import { useState, useEffect } from "react";
import { API_URL } from "../util/constants";

interface DataResponse {
  success: string;
  token: string;
}

interface AuthCheckResponse {
  code: string;
  message: string;
  data: DataResponse;
}

function useAuthCheck() {
  const [status, setStatus] = useState<"success" | "fail" | "pending">(
    "pending"
  );

  useEffect(() => {
    const checkAuth = async () => {
      const f = await fetch(`${API_URL}/auth/check`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        cache: "no-cache",
        credentials: "include",
      });
      const fJson: AuthCheckResponse = await f.json();
      if (fJson.data.success != "true" || fJson.message != "fail")
        return setStatus("fail");
      else return setStatus("success");
    };
    checkAuth();
  }, []);

  return status;
}

export default useAuthCheck;
