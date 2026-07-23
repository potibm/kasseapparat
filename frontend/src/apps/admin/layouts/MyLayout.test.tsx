import { render, screen } from "@testing-library/react";
import { AdminContext } from "react-admin";
import { MyLayout, MyAppBar } from "./MyLayout";
import { expect, test, describe, vi } from "vitest";
import { ConfigContext } from "@core/config/context/ConfigContext";
import { AppConfig } from "@core/config/types/config.types";

vi.mock("./Menu", () => ({
  Menu: () => <div data-testid="mock-menu">Mocked Menu</div>,
}));

const mockConfig: AppConfig = {
  version: "1.0.0",
  apiHost: "http://localhost:3001",
  apiBaseUrl: "http://localhost:3001/api/v3",
  websocketHost: "ws://localhost:3001",
  websocketBaseUrl: "ws://localhost:3001/api/v3",
  locale: "en",
  currencyCode: "USD",
  currencyLocale: "en-US",
  currency: new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  }),
  currencyOptions: {},
  dateLocale: "en-US",
  dateOptions: {},
  vatRates: [],
  paymentMethods: [],
  sumupEnabled: false,
  authMode: "oidc",
};

describe("Layout & AppBar", () => {
  test("MyAppBar renders title and logo", () => {
    render(
      <ConfigContext value={mockConfig}>
        <AdminContext>
          <MyAppBar />
        </AdminContext>
      </ConfigContext>,
    );

    expect(screen.getByText(/Kasseapparat/i)).toBeInTheDocument();

    const appBar = screen.getByRole("banner");
    expect(appBar).toHaveClass("MuiAppBar-colorSecondary");
  });

  test("MyLayout renders and displays the correct content", () => {
    render(
      <ConfigContext value={mockConfig}>
        <AdminContext>
          <MyLayout>
            <div data-testid="dummy-content">Content</div>
          </MyLayout>
        </AdminContext>
      </ConfigContext>,
    );

    expect(screen.getByText(/Kasseapparat/i)).toBeInTheDocument();

    expect(screen.getByTestId("dummy-content")).toBeInTheDocument();
  });

  test("Snapshot Match", () => {
    const { asFragment } = render(
      <ConfigContext value={mockConfig}>
        <AdminContext>
          <MyLayout>
            <div data-testid="dummy-content">Content</div>
          </MyLayout>
        </AdminContext>
      </ConfigContext>,
    );
    expect(asFragment()).toMatchSnapshot();
  });
});
