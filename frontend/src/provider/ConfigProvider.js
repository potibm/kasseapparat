import React, { createContext, useContext, useState, useEffect } from "react";
import PropTypes from "prop-types";

const API_HOST = process.env.REACT_APP_API_HOST ?? "http://localhost:3001";

export const ConfigContext = createContext({});

const ConfigProvider = ({ children }) => {
  const [config, setConfig] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetch(`${API_HOST}/api/v1/config`)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
      })
      .then((data) => {
        data.currencyOptions = {
          style: "currency",
          currency: data.currencyCode ?? "DKK",
          minimumFractionDigits: data.fractionDigitsMin ?? 0,
          maximumFractionDigits: data.fractionDigitsMax ?? 2,
        };

        data.currency = new Intl.NumberFormat(
          data.currencyLocale ?? "dk-DK",
          data.currencyOptions,
        );

        data.dateLocale = data.dateLocale ?? "en-US";
        try {
          data.dateOptions = JSON.parse(data.dateOptions ?? "{}");
        } catch (error) {
          console.error("Error parsing dateOptions:", error);
          data.dateOptions = {};
        }

        setConfig(data);
        setLoading(false);
      })
      .catch((error) => {
        console.error("Error fetching config:", error);
        setError(error.message);
        setLoading(false);
      });
  }, []);

  if (loading) {
    return <div>Loading Config...</div>;
  }

  if (error) {
    return <div>Error loading config: {error}</div>;
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
