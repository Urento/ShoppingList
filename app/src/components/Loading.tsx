import { Sidebar } from "./Sidebar";

interface Props {
  withSidebar?: boolean;
}

export const Loading: React.FC<Props> = ({ withSidebar = false }) => {
  return (
    <div>
      {withSidebar && <Sidebar />}
      <div className="center">
        <svg
          className="loading-page-svg"
          viewBox="25 25 50 50"
          preserveAspectRatio="xMidYMin"
        >
          <circle
            className="loading-page-circle"
            cx="50"
            cy="50"
            r="20"
          ></circle>
        </svg>
      </div>
    </div>
  );
};
