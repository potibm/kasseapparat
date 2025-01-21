import React, { createContext, useContext, useState, useEffect } from "react";
import PropTypes from "prop-types";
import MinimalBaseCard from "../components/MinimalBaseCard";
import { Alert, Spinner } from "flowbite-react";
import { HiInformationCircle } from "react-icons/hi";

const API_HOST = process.env.REACT_APP_API_HOST ?? "http://localhost:3001";

export const ConfigContext = createContext({});

const ConfigProvider = ({ children }) => {
  const [config, setConfig] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetch(`${API_HOST}/api/v2/config`)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
      })
      .then((data) => {
        data.apiHost = API_HOST;

        data.currencyOptions = {
          style: "currency",
          currency: data.currencyCode ?? "DKK",
          minimumFractionDigits: data.fractionDigitsMin ?? 0,
          maximumFractionDigits: data.fractionDigitsMax ?? 2,
        };

        data.locale = data.currencyLocale ?? "dk-DK";
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
    return (
      <MinimalBaseCard>
        Loading Config...
        <br />
        <Spinner className="mt-3" />
      </MinimalBaseCard>
    );
  }

  if (error) {
    return (
      <MinimalBaseCard>
        Error loading config:
        <Alert
          color="failure"
          icon={HiInformationCircle}
          className="mt-3 font-mono"
        >
          {error}
        </Alert>
      </MinimalBaseCard>
    );
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
