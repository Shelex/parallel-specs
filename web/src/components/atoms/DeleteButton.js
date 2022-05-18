import { useCallback } from "react";

import Spinner from "./Spinner";

export const DeleteButton = ({ onClick, loading, title }) => {
  const onConfirm = useCallback(
    (e) => {
      e.preventDefault();
      // eslint-disable-next-line no-restricted-globals
      if (confirm("Are you really sure?")) {
        return onClick(e);
      }
    },
    [onClick]
  );

  return (
    <div className="mt-10">
      <button
        className={`bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-2 rounded w-48`}
        onClick={onConfirm}
      >
        {loading ? <Spinner /> : <p>{title}</p>}
      </button>
    </div>
  );
};
