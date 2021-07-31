export const Loading: React.FC = ({}) => {
  return (
    <div className="center">
      <svg
        className="loading-page-svg"
        viewBox="25 25 50 50"
        preserveAspectRatio="xMidYMin"
      >
        <circle className="loading-page-circle" cx="50" cy="50" r="20"></circle>
      </svg>
    </div>
  );
};
