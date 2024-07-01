import React from "react";
import { Card } from "flowbite-react";
import { useConfig } from "../provider/ConfigProvider";
import PropTypes from "prop-types";
import { Link } from "react-router-dom";

const BaseCard = ({
  children,
  title = null,
  linkLogin = false,
  linkForgotPassword = false,
}) => {
  const version = useConfig().version;

  return (
    <div className="flex justify-center items-center h-screen">
      <Card className="max-w-sm ">
        <h5 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
          <Link to="/">
            <img
              src="/android-chrome-192x192.png"
              alt="Kasseapparat"
              className="align-text-top h-7 inline"
            />{" "}
            Kasseapparat
          </Link>
        </h5>

        <div className="my-3">
          {title && <h2 className="text-xl mb-2">{title}</h2>}

          {children}
        </div>

        <hr />
        <p className="text-xs">
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
          Version {version}
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
