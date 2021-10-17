import { useHistory } from "react-router";
import { Loading } from "../components/Loading";
import { Sidebar } from "../components/Sidebar";
import useAuthCheck from "../hooks/useAuthCheck";

const NotFound: React.FC = ({}) => {
  const authStatus = useAuthCheck();
  const history = useHistory();
  if (authStatus === "fail") {
    localStorage.removeItem("authenticated");
    history.push("/");
  }

  if (authStatus === "pending") return <Loading withSidebar />;

  return (
    <div className="flex flex-no-wrap h-screen">
      <Sidebar />
      <div className="container mx-auto py-10 md:w-4/5 w-11/12">
        <div className="w-full h-full rounded">
          <div className="flex items-center justify-between">
            <h1>Page not found</h1>
          </div>
        </div>
      </div>
    </div>
  );
};

export default NotFound;
