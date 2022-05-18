import { memo, useCallback } from "react";
import { BiExit } from "react-icons/bi";
import { withRouter } from "react-router-dom";
import { auth } from "../../services/auth.service";

const Logout = ({ history, className }) => {
  const onClick = useCallback(() => {
    auth.logout();
    history.push("/");
  }, [history]);

  return (
    <button
      className={`focus:outline-none hover:bg-blue-900 p-1 rounded-md ${className}`}
      onClick={onClick}
    >
      <BiExit className="text-white text-2xl sm:text-3xl" />
    </button>
  );
};

export default withRouter(memo(Logout));
