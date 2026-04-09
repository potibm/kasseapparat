import React from "react"; // WICHTIG: React importieren für React.useState
import {
  render,
  screen,
  waitFor,
  act,
  fireEvent,
} from "@testing-library/react"; // WICHTIG: fireEvent hinzugefügt
import { describe, expect, it, vi, beforeEach } from "vitest";
import { AuthProvider } from "./AuthProvider";
import { useAuth } from "../hooks/useAuth";
import * as AuthApi from "@core/api/auth";
import * as AuthStorage from "../services/auth-storage";
import * as ConfigHook from "@core/config/hooks/useConfig";
import { RefreshTokenResponse } from "@core/api/auth.schemas";

// Mocks
vi.mock("@core/api/auth");
vi.mock("../services/auth-storage");
vi.mock("@core/config/hooks/useConfig");

const TestConsumer = () => {
  const { username, getToken } = useAuth();
  const [token, setToken] = React.useState<string | null>(null);

  const handleGetToken = async () => {
    const t = await getToken();
    setToken(t);
  };

  return (
    <div>
      <span data-testid="username">{username}</span>
      <span data-testid="token">{token}</span>
      <button onClick={handleGetToken}>Get Token</button>
    </div>
  );
};

describe("AuthProvider", () => {
  beforeEach(() => {
    vi.resetAllMocks();
    vi.mocked(ConfigHook.useConfig).mockReturnValue({
      apiHost: "http://api.test",
    } as unknown as ReturnType<typeof ConfigHook.useConfig>);

    vi.mocked(AuthStorage.getInitialUser).mockReturnValue(null);
    vi.mocked(AuthStorage.getInitialSession).mockReturnValue({
      token: null,
      expiryDate: null,
    });
  });

  it("provides initial user data from storage", () => {
    // 1. Setup
    vi.mocked(AuthStorage.getInitialUser).mockReturnValue({
      id: 123,
      username: "Stefan",
      role: "admin",
      gravatarUrl: "",
    });

    // 2. Render
    render(
      <AuthProvider>
        <TestConsumer />
      </AuthProvider>,
    );

    // 3. Assert
    expect(screen.getByTestId("username")).toHaveTextContent("Stefan");
  });

  it("handles token refresh when token is expired", async () => {
    // 1. Setup: expired token in storage
    const expiredDate = new Date(Date.now() - 10000); // 10s in der Vergangenheit
    vi.mocked(AuthStorage.getInitialSession).mockReturnValue({
      token: "old-token",
      expiryDate: expiredDate,
    });

    // API Mock
    vi.mocked(AuthApi.refreshJwtToken).mockResolvedValue({
      access_token: "new-token",
      expires_in: 3600,
    } as unknown as RefreshTokenResponse);

    render(
      <AuthProvider>
        <TestConsumer />
      </AuthProvider>,
    );

    // 2. Action
    const btn = screen.getByText("Get Token");
    await act(async () => {
      fireEvent.click(btn);
    });

    // 3. Assert
    await waitFor(() => {
      expect(AuthApi.refreshJwtToken).toHaveBeenCalledWith("http://api.test");
      expect(screen.getByTestId("token")).toHaveTextContent("new-token");
    });

    expect(AuthStorage.storeSession).toHaveBeenCalledWith(
      "new-token",
      expect.any(Date),
    );
  });

  it("clears session on refresh failure", async () => {
    // 1. Setup: Expired token
    const expiredDate = new Date(Date.now() - 10000);
    vi.mocked(AuthStorage.getInitialSession).mockReturnValue({
      token: "old-token",
      expiryDate: expiredDate,
    });

    // API fails
    vi.mocked(AuthApi.refreshJwtToken).mockRejectedValue(
      new Error("Unauthorized"),
    );

    render(
      <AuthProvider>
        <TestConsumer />
      </AuthProvider>,
    );

    // 2. Action
    const btn = screen.getByText("Get Token");
    await act(async () => {
      fireEvent.click(btn);
    });

    // 3. Assert
    await waitFor(() => {
      expect(AuthStorage.clearAuthStorage).toHaveBeenCalled();
      expect(screen.getByTestId("token")).toBeEmptyDOMElement();
    });
  });

  it("prevents multiple simultaneous refresh API calls (Race Condition)", async () => {
    // 1. Setup: Expired token
    const expiredDate = new Date(Date.now() - 10000);
    vi.mocked(AuthStorage.getInitialSession).mockReturnValue({
      token: "old-token",
      expiryDate: expiredDate,
    });

    // we delay the API response artificially so that we can click again "in the meantime"
    vi.mocked(AuthApi.refreshJwtToken).mockImplementation(
      () =>
        new Promise((resolve) =>
          setTimeout(
            () =>
              resolve({
                access_token: "new-token",
                expires_in: 3600,
              } as unknown as RefreshTokenResponse),
            100,
          ),
        ),
    );

    render(
      <AuthProvider>
        <TestConsumer />
      </AuthProvider>,
    );

    const btn = screen.getByText("Get Token");

    // 2. Action: we click twice in quick succession
    await act(async () => {
      fireEvent.click(btn);
      fireEvent.click(btn);
    });

    // 3. Assert: Wait, until the token has arrived
    await waitFor(() => {
      expect(screen.getByTestId("token")).toHaveTextContent("new-token");
    });

    // The most important: The API should have been called only once!
    expect(AuthApi.refreshJwtToken).toHaveBeenCalledTimes(1);
  });
});
