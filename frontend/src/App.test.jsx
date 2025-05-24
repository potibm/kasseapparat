import { render, screen } from "@testing-library/react";
import App from "./App";
import PropTypes from "prop-types";
import AuthProvider from "./Auth/provider/AuthProvider";
import { ConfigContext } from "./provider/ConfigProvider";
import Login from "./Auth/components/Login";
import { createRoutesStub } from "react-router";
import { describe, expect, it } from "vitest";

import React from "react";

const mockConfig = {
  version: "0.2.0",
};

const MockConfigProvider = ({ children }) => (
  <ConfigContext.Provider value={mockConfig}>{children}</ConfigContext.Provider>
);

MockConfigProvider.propTypes = {
  children: PropTypes.node.isRequired,
};

describe("App", () => {
  it("renders the App component before loading the config", () => {
    render(<App />);
    const headlineElement = screen.getByText(/Loading Config/i);
    expect(headlineElement).toBeInTheDocument();
  });

  it("renders the Login component", () => {
    const LoginComponentWrapped = () => {
      return (
        <div>
          <MockConfigProvider>
            <AuthProvider>
              <Login />
            </AuthProvider>
          </MockConfigProvider>
        </div>
      );
    };

    const Stub = createRoutesStub([
      {
        index: true,
        path: "/",
        Component: LoginComponentWrapped,
      },
    ]);

    // render the app stub at "/login"
    render(<Stub initialEntries={["/"]} />);

    const headlineElement = screen.getByText(/Kasseapparat/i);
    expect(headlineElement).toBeInTheDocument();

    const versionElement = screen.getByText(/Version 0.2.0/i);
    expect(versionElement).toBeInTheDocument();
  });
});
