import React from "react";
import { Card, DarkThemeToggle } from "flowbite-react";
import Version from "./Version";
import PropTypes from "prop-types";
import { Link } from "react-router";

const BaseCard = ({
  children,
  title = null,
  linkLogin = false,
  linkForgotPassword = false,
}) => {
  return (
    <div className="flex justify-center items-center h-screen">
      <Card className="max-w-sm ">
        <h5 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-gray-200">
          <Link to="/">
            <img
              src="/android-chrome-192x192.png"
              alt="Kasseapparat"
              className="align-text-top h-7 inline"
            />{" "}
            Kasseapparat
          </Link>
        </h5>

        <div className="my-3 dark:text-gray-200">
          {title && <h2 className="text-xl mb-2 dark:text-white">{title}</h2>}

          {children}
        </div>

        <hr />
        <p className="text-xs dark:text-gray-200">
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
        </p>
      </Card>
    </div>
  );
};

BaseCard.propTypes = {
  children: PropTypes.node.isRequired,
  title: PropTypes.string,
  linkLogin: PropTypes.bool,
  linkForgotPassword: PropTypes.bool,
};

export default BaseCard;
