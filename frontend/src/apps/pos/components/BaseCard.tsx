import React, { ReactNode } from "react";
import { DarkThemeToggle } from "flowbite-react";
import Version from "./Version";
import { Link } from "react-router";
import MinimalBaseCard from "./MinimalBaseCard";

interface BaseCardProps {
  children: ReactNode;
  title: string;
  linkLogin?: boolean;
  linkForgotPassword?: boolean;
}

export const BaseCard: React.FC<BaseCardProps> = ({
  children,
  title = "",
  linkLogin = false,
  linkForgotPassword = false,
}) => {
  return (
    <MinimalBaseCard
      title={title}
      navigation={
        <>
          {linkLogin && (
            <>
              <Link to="/" className="text-blue-500 hover:underline">
                Login
              </Link>
              <span className="mx-2">&ndash;</span>
            </>
          )}
          {linkForgotPassword && (
            <>
              <Link
                to="/forgot-password"
                className="text-blue-500 hover:underline"
              >
                Forgot password
              </Link>
              <span className="mx-2">&ndash;</span>
            </>
          )}
          <Link
            to="/manual.pdf"
            reloadDocument
            className="text-blue-500 hover:underline"
          >
            Manual
          </Link>
          <span className="mx-2">&ndash;</span>
          <Version />
          <DarkThemeToggle aria-label="Toggle dark mode" className="mx-2" />
        </>
      }
    >
      {children}
    </MinimalBaseCard>
  );
};

export default BaseCard;
