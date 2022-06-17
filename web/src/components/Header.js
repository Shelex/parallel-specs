import { useState, useCallback, useEffect } from "react";
import { Link } from "react-router-dom";

import { Logo } from "./atoms/Logo";
import Logout from "./atoms/Logout";
import Menu from "./atoms/Menu";
import { Nav } from "./atoms/Nav";

import { auth, authChangedEvent } from "../services/auth.service";

const menuItems = [
  { name: "Emulate Session", link: "emulate" },
  {
    name: "API Keys",
    link: "apiKeys",
  },
];

export const Header = ({ title }) => {
  const [isMenu, setIsMenu] = useState(false);

  const parseUsername = () => (auth.get() ? auth.info().email : "");

  const [username, setUsername] = useState(parseUsername());

  useEffect(() => {
    const usernameHandler = () => setUsername(parseUsername());
    window.addEventListener(authChangedEvent, usernameHandler, false);
  }, []);

  const onClick = useCallback(() => {
    setIsMenu((prev) => !prev);
  }, []);

  return (
    <header className="bg-blue-800 max-w-full">
      <div className="h-14 sm:h-16 max-w-2xl mx-auto px-4 flex items-center justify-between">
        <div className="flex-1 flex items-center justify-between sm:justify-start">
          <Menu isMenu={isMenu} onClick={onClick} />
          <Logo title={title} />
          <Nav items={menuItems} />
          <Logout className="block sm:hidden" />
        </div>
        <p className="text-white px-2">{username}</p>
        <Logout className="hidden sm:block" />
      </div>
      {isMenu && (
        <div className="sm:hidden px-4 py-2 pb-6 space-y-1">
          {menuItems.map((item, index) => (
            <Link key={index} className="menu-link" to={`/${item.link}`}>
              {item.name}
            </Link>
          ))}
        </div>
      )}
    </header>
  );
};
