import { useState } from "react";
import { useEffect } from "react";
import { BackupCodesResponse } from "../../../types/BackupCodes";
import { API_URL } from "../../../util/constants";

type BackupCodeState = {
  codes: string[];
};

const useLoadBackupCodes = (condition: boolean) => {
  const [loadingBackupCodes, setLoadingBackupCodes] = useState<boolean>(true);
  const [backupCodes, setBackupCodes] = useState<BackupCodeState>({
    codes: [],
  });
  const [has, setHas] = useState<boolean>(false);

  const loadBackupCodes = async () => {
    const response = await fetch(`${API_URL}/backupcodes`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    });
    const fJson: BackupCodesResponse = await response.json();
    if (fJson.data.has && fJson.data.has === "false") {
      setHas(false);
    } else if (fJson.data.success === "true") {
      const c = fJson.data.codes.replace("{", "").replace("}", "");
      const codesToArray = c.split(",");
      console.log(codesToArray);
      setBackupCodes({ codes: codesToArray });
      setHas(true);
    } else {
      setHas(false);
    }
    console.log(fJson.data.codes);
    setLoadingBackupCodes(false);
  };

  useEffect(() => {
    loadBackupCodes();
  }, []);

  useEffect(() => {
    if (condition) loadBackupCodes();
  }, [condition]);

  return {
    loadingBackupCodes,
    backupCodes,
    setBackupCodes,
    setLoadingBackupCodes,
    has,
  };
};

export default useLoadBackupCodes;
