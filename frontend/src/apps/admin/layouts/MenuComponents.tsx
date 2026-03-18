import React from "react";

interface MenuDividerProps {
  name: string;
}

export const MenuDivider: React.FC<MenuDividerProps> = ({ name }) => (
  <div className="text-right mt-4 text-xs pr-1 font-bold tracking-wide border-b-2 border-b-pink-500 uppercase text-ellipsis">
    {name}
  </div>
);

interface MenuMessageProps {
  message: string;
}

export const MenuMessage: React.FC<MenuMessageProps> = ({ message }) => (
  <div className="text-left mt-4 text-xs font-bold tracking-wide border-t-2 border-b-2 border-t-red-700 border-b-red-700 bg-info pl-5 pr-3 uppercase text-ellipsis">
    {message}
  </div>
);
