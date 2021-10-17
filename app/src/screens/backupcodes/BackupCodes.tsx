import { useState } from "react";
import { useHistory } from "react-router";
import swal from "sweetalert";
import { Button } from "../../components/Button";
import { Loading } from "../../components/Loading";
import { Sidebar } from "../../components/Sidebar";
import useAuthCheck from "../../hooks/useAuthCheck";
import { GenerateBackupCodesResponse } from "../../types/BackupCodes";
import { API_URL } from "../../util/constants";
import useLoadBackupCodes from "./hooks/useLoadBackupCodes";

const BackupCodes: React.FC = () => {
  const authStatus = useAuthCheck();
  const history = useHistory();
  const [refreshCodes, setRefreshCodes] = useState<boolean>(false);
  const {
    loadingBackupCodes,
    backupCodes,
    setBackupCodes,
    setLoadingBackupCodes,
    has,
  } = useLoadBackupCodes(refreshCodes);

  if (loadingBackupCodes) return <Loading withSidebar />;
  if (authStatus === "fail") {
    localStorage.removeItem("authenticated");
    history.push("/");
  }

  if (authStatus === "pending") return <Loading withSidebar />;

  const generateBackupCodes = async (regenerate = false) => {
    setLoadingBackupCodes(true);
    const response = await fetch(
      `${API_URL}/backupcodes/${regenerate ? "regenerate" : "generate"}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        credentials: "include",
      }
    );
    const fJson: GenerateBackupCodesResponse = await response.json();
    if (fJson.data.success === "true") {
      const c = fJson.data.codes.replace("{", "").replace("}", "");
      const codesToArray = c.split(",");
      setBackupCodes({ codes: codesToArray });
      if (regenerate) setRefreshCodes(true);
    }
    setLoadingBackupCodes(false);
  };

  const regenerateBackupCodes = async () => {
    const alert = await swal({
      icon: "warning",
      title: "Are you sure you want to regenerate your Backup Codes?",
      buttons: ["No, don't regenerate!", "Yes, regenerate!"],
    });

    if (alert) {
      swal.close!();
      swal({
        title: "Regenerating...",
      });
      await generateBackupCodes(true);
      swal.close!();
      swal({
        icon: "success",
        title: "Successfully regenerated your Backup Codes",
      });
    }
  };

  backupCodes.codes.map((e: string) => {
    console.log(e);
  });

  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
      <div className="container mx-auto py-10 md:w-4/5 w-11/12">
        <div className="bg-white px-4 md:px-10 pt-5 md:pt-7 pb-5 overflow-y-auto">
          <h1 className="text-4xl text-center">Backup Codes</h1>
          {backupCodes.codes.length <= 0 && !loadingBackupCodes && (
            <Button
              text="Generate Backup Codes"
              loadingText="Generating Backup Codes..."
              onClick={() => generateBackupCodes(false)}
              className="inline-flex sm:ml-3 mt-4 sm:mt-0 items-start justify-start px-6 py-3 bg-green-600 hover:bg-green-500 text-white focus:outline-none rounded"
              loading={loadingBackupCodes}
              color="green"
            />
          )}
          {has && !loadingBackupCodes && (
            <Button
              text="Regenerate Backup Codes"
              loadingText="Regenerating Backup Codes..."
              onClick={regenerateBackupCodes}
              className="inline-flex sm:ml-3 mt-4 sm:mt-0 items-start justify-start px-6 py-3 bg-green-600 hover:bg-green-500 text-white focus:outline-none rounded"
              loading={loadingBackupCodes}
              color="green"
            />
          )}
          <table className="w-full whitespace-nowrap">
            <thead>
              <tr className="h-16 w-full text-sm leading-none ">
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
                <th className="font-normal text-left pl-12"></th>
              </tr>
            </thead>
            <tbody className="w-full">
              {backupCodes.codes.length > 0 && (
                <tr className="h-20 text-lg leading-none text-gray-800 bg-white hover:bg-gray-100 border-b border-t border-gray-100">
                  <td className="pl-12">
                    <p className="font-medium">{backupCodes.codes[0]}</p>
                  </td>
                  <td className="pl-12">
                    <p className="font-medium">{backupCodes.codes[1]}</p>
                  </td>
                  <td className="pl-12">
                    <p className="font-medium">{backupCodes.codes[2]}</p>
                  </td>
                  <td className="pl-12">
                    <p className="font-medium">{backupCodes.codes[3]}</p>
                  </td>
                  <td className="pl-12">
                    <p className="font-medium">{backupCodes.codes[4]}</p>
                  </td>
                  <td className="pl-12">
                    <p className="font-medium">{backupCodes.codes[5]}</p>
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default BackupCodes;
