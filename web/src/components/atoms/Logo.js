import { Link } from "react-router-dom";
import { BiServer } from "react-icons/bi";

export const Logo = ({ title }) => {
  return (
    <Link to="/">
      <h1 className="flex items-center hover:bg-blue-900 px-2 py-1 rounded-md text-white text-xl sm:text-2xl font-semibold sm:pl-1">
        <BiServer className="hidden sm:block text-3xl" />
        {title}
      </h1>
    </Link>
  );
};
