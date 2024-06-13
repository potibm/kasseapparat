import React from "react";
import { render, screen } from "@testing-library/react";
import App from "./App";
import { test, expect } from "@jest/globals";
import Routes from "./routes";
import AuthProvider from "./provider/AuthProvider";

test("renders the application before loading the config", () => {
  render(<App />);
  const headlineElement = screen.getByText(/Loading Config/i);
  expect(headlineElement).toBeInTheDocument();
});

test("renders the application after loading the config", () => {
  render(
    <AuthProvider>
      <Routes />
    </AuthProvider>,
  );
  const headlineElement = screen.getByText(/Kasseapparat/i);
  expect(headlineElement).toBeInTheDocument();
});
