import React from "react";
import { render, screen } from "@testing-library/react";
import App from "./App";
import { test, expect } from "@jest/globals";
import Routes from "./routes";
import AuthProvider from "./provider/AuthProvider";
import { ConfigContext } from "./provider/ConfigProvider";
import PropTypes from "prop-types";

const mockConfig = {
  version: "1.0.0",
};

const MockConfigProvider = ({ children }) => (
  <ConfigContext.Provider value={mockConfig}>{children}</ConfigContext.Provider>
);

MockConfigProvider.propTypes = {
  children: PropTypes.node.isRequired,
};

test("renders the application before loading the config", () => {
  render(<App />);
  const headlineElement = screen.getByText(/Loading Config/i);
  expect(headlineElement).toBeInTheDocument();
});

test("renders the application after loading the config", () => {
  render(
    <MockConfigProvider>
      <AuthProvider>
        <Routes />
      </AuthProvider>
    </MockConfigProvider>,
  );
  const headlineElement = screen.getByText(/Kasseapparat/i);
  expect(headlineElement).toBeInTheDocument();

  const versionElement = screen.getByText(/Version 1.0.0/i);
  expect(versionElement).toBeInTheDocument();
});
