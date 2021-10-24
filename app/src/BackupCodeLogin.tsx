import { useHistory } from "react-router";
import { Button } from "./components/Button";
import { useState, useEffect } from "react";
import { API_URL } from "./util/constants";
import { isLoggedIn } from "./storage/UserStorage";
import { ResetPasswordBackupCode, VerifyBackupCode } from "./types/BackupCodes";
import swal from "sweetalert";

const BackupCodeLogin: React.FC = () => {
  const history = useHistory();
  const loggedIn = isLoggedIn();
  const [backupCode, setBackupCode] = useState<string>("");
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [repeatPassword, setRepeatPassword] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [canReset, setCanReset] = useState({
    canReset: false,
    code: "",
    owner: "",
  });

  useEffect(() => {
    //redirect for checking if the user is actually already authenticated
    if (loggedIn) history.push("/dashboard");
  }, []);

  const handleBackupCodeChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setBackupCode(e.target.value);
  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setEmail(e.target.value);
  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setPassword(e.target.value);
  const handleRepeatPasswordChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setRepeatPassword(e.target.value);

  const verifyBackupCode = async (event: any) => {
    event.preventDefault();
    setLoading(true);
    setCanReset({ canReset: false, code: "", owner: "" });

    const response = await fetch(`${API_URL}/backupcodes`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify({
        owner: email,
        code: backupCode,
      }),
    });
    const fJson: VerifyBackupCode = await response.json();
    if (fJson) setLoading(false);

    if (
      !fJson.data ||
      fJson.data.error ||
      fJson.data.success === "false" ||
      fJson.data.ok === "false"
    ) {
      swal({
        icon: "error",
        title: "Backup Code is incorrect",
      });
      return;
    }
    setCanReset({
      canReset: true,
      code: backupCode,
      owner: email,
    });
    swal({
      icon: "success",
      title: "Backup Code is correct",
      text: "Reset your password now!",
    });
  };

  const resetPassword = async (event: any) => {
    event.preventDefault();
    setLoading(true);

    if (password !== repeatPassword) {
      setLoading(false);
      swal({
        icon: "error",
        title: "Both Passwords have to be the same",
      });
      return;
    }

    const response = await fetch(`${API_URL}/backupcodes/changepassword`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      cache: "no-cache",
      credentials: "include",
      body: JSON.stringify({
        code: canReset.code,
        owner: canReset.owner,
        password: password,
      }),
    });
    const fJson: ResetPasswordBackupCode = await response.json();
    if (fJson) setLoading(false);
    if (
      !fJson.data ||
      fJson.data.error ||
      fJson.data.success === "false" ||
      fJson.data.ok === "false"
    ) {
      swal({
        icon: "error",
        title: "Error while resetting password! Try again later!",
      });
      return;
    }
    setCanReset({
      canReset: true,
      code: backupCode,
      owner: email,
    });
    swal({
      icon: "success",
      title: "Successfully reset your password!",
      text: "You will be redirected shortly!",
    });
    setTimeout(() => history.push("/"), 2000);
  };

  if (!canReset.canReset) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-md w-full space-y-8">
          <div>
            <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
              Backup Code Login
            </h2>
          </div>
          <form className="mt-8 space-y-6" onSubmit={verifyBackupCode}>
            <input type="hidden" name="remember" defaultValue="true" />
            <div className="rounded-md shadow-sm -space-y-px">
              <div>
                <label htmlFor="email-address" className="sr-only">
                  Email
                </label>
                <input
                  id="email-address"
                  name="email-address"
                  type="email"
                  autoComplete="email"
                  required
                  className={`text-sm sm:text-base relative w-full border rounded placeholder-gray-400 focus:border-indigo-400 focus:outline-none py-2 pr-2`}
                  placeholder="Email"
                  onChange={handleEmailChange}
                />
              </div>

              <div style={{ paddingTop: "2%" }}>
                <label htmlFor="backupcode" className="sr-only">
                  Backup Code
                </label>
                <input
                  id="backupcode"
                  name="backupcode"
                  type="text"
                  autoComplete="backupcode"
                  required
                  className={`text-sm sm:text-base relative w-full border rounded placeholder-gray-400 focus:border-indigo-400 focus:outline-none py-2 pr-2`}
                  placeholder="Backup Code"
                  onChange={handleBackupCodeChange}
                />
              </div>
            </div>

            <div className="flex items-center justify-between">
              <div className="text-sm">
                <a
                  onClick={() => history.push("/")}
                  href="#"
                  className="font-medium text-indigo-600 hover:text-indigo-500"
                >
                  Normal Login
                </a>
              </div>
            </div>

            <div>
              <Button
                showIcon={true}
                loading={loading}
                onClick={verifyBackupCode}
                text="Verify"
                loadingText="Verifying..."
                type="submit"
              ></Button>
            </div>
          </form>
        </div>
      </div>
    );
  } else {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-md w-full space-y-8">
          <div>
            <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
              Backup Code Login <br /> Reset Password
            </h2>
          </div>
          <form className="mt-8 space-y-6" onSubmit={resetPassword}>
            <input type="hidden" name="remember" defaultValue="true" />
            <div className="rounded-md shadow-sm -space-y-px">
              <div>
                <label htmlFor="password" className="sr-only">
                  New Password
                </label>
                <input
                  id="password"
                  name="password"
                  type="password"
                  autoComplete="password"
                  required
                  className={`text-sm sm:text-base relative w-full border rounded placeholder-gray-400 focus:border-indigo-400 focus:outline-none py-2 pr-2`}
                  placeholder="New Password"
                  onChange={handlePasswordChange}
                />
              </div>

              <div style={{ paddingTop: "2%" }}>
                <label htmlFor="repeat-password" className="sr-only">
                  Repeat Password
                </label>
                <input
                  id="repeat-password"
                  name="repeat-password"
                  type="password"
                  autoComplete="password"
                  required
                  className={`text-sm sm:text-base relative w-full border rounded placeholder-gray-400 focus:border-indigo-400 focus:outline-none py-2 pr-2`}
                  placeholder="Repeat Password"
                  onChange={handleRepeatPasswordChange}
                />
              </div>
            </div>

            <div className="flex items-center justify-between">
              <div className="text-sm">
                <a
                  onClick={() => history.push("/")}
                  href="#"
                  className="font-medium text-indigo-600 hover:text-indigo-500"
                >
                  Normal Login
                </a>
              </div>
            </div>

            <div>
              <Button
                showIcon={true}
                loading={loading}
                onClick={resetPassword}
                text="Reset Password"
                loadingText="Resetting Password..."
                type="submit"
              ></Button>
            </div>
          </form>
        </div>
      </div>
    );
  }
};

export default BackupCodeLogin;
