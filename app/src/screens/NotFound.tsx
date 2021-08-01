import { Sidebar } from "../components/Sidebar";

export const NotFound: React.FC = ({}) => {
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
