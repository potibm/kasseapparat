import { render, screen, waitFor } from "@testing-library/react";
import { createRoutesStub } from "react-router";
import { describe, expect, it, vi, beforeEach, afterEach } from "vitest";
import * as ConfigProviderModule from "./core/config/providers/ConfigProvider";
import AuthProvider from "./apps/pos/features/auth/providers/auth-provider";
import Login from "./apps/pos/features/auth/components/Login";
import App from "./App";
import "@testing-library/jest-dom";
import { ReactNode } from "react";

vi.stubGlobal(
  "fetch",
  vi.fn().mockResolvedValue({
    ok: true,
    json: () => Promise.resolve({}),
  }),
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
    vi.unstubAllGlobals(); // Extrem wichtig für stubGlobal('fetch')!
  });

  it("renders the App component before loading the config", async () => {
    const spy = vi.spyOn(ConfigProviderModule, "useConfig").mockReturnValue({
      config: null,
      loading: true,
      error: null,
    } as any);

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
      config: { version: "0.2.0" },
      loading: false,
      error: null,
      version: "0.2.0",
      apiHost: "http://localhost",
      currency: { format: (v: number) => `${v} €` },
      paymentMethods: [],
    };

    vi.spyOn(ConfigProviderModule, "useConfig").mockReturnValue(
      mockConfigValue as any,
    );

    const MockConfigProvider = ({ children }: { children: ReactNode }) => (
      <ConfigProviderModule.ConfigContext.Provider
        value={mockConfigValue as any}
      >
        {children}
      </ConfigProviderModule.ConfigContext.Provider>
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
