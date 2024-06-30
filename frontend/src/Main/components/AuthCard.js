import React from "react";
import { Card } from "flowbite-react";
import { useConfig } from "../../provider/ConfigProvider";
import PropTypes from "prop-types";
import { Link } from "react-router-dom";

const AuthCard = ({ children, title = null }) => {
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
        <p className="text-xs">Version {version}</p>
      </Card>
    </div>
  );
};

AuthCard.propTypes = {
  children: PropTypes.node.isRequired,
  title: PropTypes.string,
};

export default AuthCard;
