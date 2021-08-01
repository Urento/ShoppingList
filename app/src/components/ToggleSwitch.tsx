import clsx from "clsx";
import { MouseEventHandler } from "react";

interface Props {
  toggled: boolean;
  onClick: MouseEventHandler<HTMLDivElement> | undefined;
  id: string;
  title: string;
}

export const ToggleSwitch: React.FC<Props> = ({
  toggled,
  onClick,
  id,
  title,
}) => {
  return (
    <div className="flex justify-between items-center" onClick={onClick}>
      <h2 className="block mb-2 text-sm font-bold text-gray-700">{title}</h2>
      <div
        className={clsx(
          `w-16 h-10 flex items-center bg-gray-300 rounded-full p-1 duration-300 ease-in-out ${
            toggled && "bg-green-400"
          }`
        )}
      >
        <div
          id={id}
          className={clsx(
            `bg-white w-8 h-8 rounded-full shadow-md transform duration-300 ease-in-out ${
              toggled && "translate-x-6"
            }`
          )}
        ></div>
      </div>
    </div>
  );
};
