import React, { createContext, useContext, useState, useEffect } from "react";
import PropTypes from "prop-types";

const API_HOST = process.env.REACT_APP_API_HOST ?? "http://localhost:3001";

const ConfigContext = createContext(null);

const ConfigProvider = ({ children }) => {
  const [config, setConfig] = useState(null);

  useEffect(() => {
    fetch(`${API_HOST}/api/v1/config`)
      .then((response) => response.json())
      .then((data) => {
        data.currencyOptions = {
          style: "currency",
          currency: data.currencyCode ?? "DKK",
          minimumFractionDigits: data.fractionDigitsMin ?? 0,
          maximumFractionDigits: data.fractionDigitsMax ?? 2,
        };

        data.currency = new Intl.NumberFormat(
          data.Locale ?? "dk-DK",
          data.currencyOptions,
        );

        setConfig(data);
      })
      .catch((error) => console.error("Error fetching config:", error));
  }, []);

  if (!config) {
    return <div>Loading Config...</div>;
  }

  return (
    <ConfigContext.Provider value={config}>{children}</ConfigContext.Provider>
  );
};

export const useConfig = () => useContext(ConfigContext);

ConfigProvider.propTypes = {
  children: PropTypes.node.isRequired,
};

export default ConfigProvider;
