import { LockClosedIcon } from "@heroicons/react/solid";
import clsx from "clsx";
import { FormEventHandler } from "react";

interface Props {
  text: string;
  onClick: FormEventHandler<HTMLButtonElement>;
  loading: boolean;
  showIcon: boolean;
  type: "button" | "reset" | "submit";
}

export const Button: React.FC<Props> = ({
  text,
  onClick,
  loading = false,
  showIcon = false,
  type = "submit",
}) => {
  return (
    <button
      type={type}
      className={clsx(
        `group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500`
      )}
      onSubmit={onClick}
    >
      {loading && (
        <svg className="loading-svg" viewBox="25 25 50 50">
          <circle className="loading-circle" cx="50" cy="50" r="20"></circle>
        </svg>
      )}
      {showIcon && !loading ? (
        <span
          className={clsx(`absolute left-0 inset-y-0 flex items-center pl-3`)}
        >
          <LockClosedIcon
            className="h-5 w-5 text-indigo-500 group-hover:text-indigo-400"
            aria-hidden="true"
          />
        </span>
      ) : (
        ""
      )}
      {!loading && text}
    </button>
  );
};
