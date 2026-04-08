import { render, screen, waitFor } from "@testing-library/react";
import { createRoutesStub } from "react-router";
import { describe, expect, it, vi, beforeEach, afterEach } from "vitest";
import * as ConfigHookModule from "./core/config/hooks/useConfig";
import { ConfigContext } from "./core/config/context/ConfigContext";
import AuthProvider from "./apps/pos/features/auth/providers/AuthProvider";
import Login from "./apps/pos/features/auth/components/Login";
import App from "./App";
import "@testing-library/jest-dom";
import { ReactNode } from "react";
import { AppConfig } from "@core/config/types/config.types";

vi.stubGlobal(
  "fetch",
  vi.fn().mockResolvedValue({
    ok: true,
    json: () => Promise.resolve({}),
  } as Response),
);

describe("App", () => {
  beforeEach(() => {
    vi.resetAllMocks();
    vi.restoreAllMocks();
  });

  afterEach(() => {
    vi.clearAllMocks();
    vi.resetAllMocks();
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
  });

  it("renders the App component before loading the config", async () => {
    const spy = vi.spyOn(ConfigHookModule, "useConfig").mockReturnValue({
      loading: true,
      error: null,
    } as unknown as AppConfig);

    render(<App />);

    await waitFor(
      () => {
        const element = screen.getByText(/loading config/i);
        expect(element).toBeInTheDocument();
      },
      { timeout: 2000 },
    );

    spy.mockRestore();
  });

  it("renders the Login component when config is loaded", async () => {
    const mockConfigValue = {
      version: "0.2.0",
      apiHost: "http://localhost",
      currency: new Intl.NumberFormat("de-DE", {
        style: "currency",
        currency: "EUR",
      }),
      paymentMethods: [],
      currencyOptions: {},
      sumupEnabled: false,
      websocketHost: "ws://localhost",
    } as unknown as AppConfig;

    vi.spyOn(ConfigHookModule, "useConfig").mockReturnValue(mockConfigValue);

    const MockConfigProvider = ({ children }: { children: ReactNode }) => (
      <ConfigContext.Provider value={mockConfigValue}>
        {children}
      </ConfigContext.Provider>
    );

    const LoginComponentWrapped = () => (
      <MockConfigProvider>
        <AuthProvider>
          <Login />
        </AuthProvider>
      </MockConfigProvider>
    );

    const Stub = createRoutesStub([
      {
        path: "/",
        Component: LoginComponentWrapped,
      },
    ]);

    render(<Stub initialEntries={["/"]} />);

    const title = await screen.findByText(/Kasseapparat/i);
    expect(title).toBeInTheDocument();

    const version = screen.getByText(/Version 0.2.0/i);
    expect(version).toBeInTheDocument();
  });
});
