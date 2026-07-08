import React from "react";
import { render, screen, act, fireEvent } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { AuthProvider } from "./AuthProvider";
import { useAuth } from "../hooks/useAuth";

// Phase 1: Tests for dummy auth bypass
const TestConsumer = () => {
  const { username, role, id, getToken, isLoggedIn } = useAuth();
  const [token, setToken] = React.useState<string | null>(null);
  const [loggedIn, setLoggedIn] = React.useState<boolean | null>(null);

  const handleGetToken = async () => {
    const t = await getToken();
    setToken(t);
  };

  const handleCheckLoggedIn = async () => {
    const status = await isLoggedIn();
    setLoggedIn(status);
  };

  return (
    <div>
      <span data-testid="username">{username}</span>
      <span data-testid="role">{role}</span>
      <span data-testid="id">{id}</span>
      <span data-testid="token">{token}</span>
      <span data-testid="loggedIn">{loggedIn?.toString() ?? "null"}</span>
      <button onClick={handleGetToken}>Get Token</button>
      <button onClick={handleCheckLoggedIn}>Check Logged In</button>
    </div>
  );
};

describe("AuthProvider (Phase 1 Dummy Bypass)", () => {
  it("provides dummy user data with admin role", () => {
    render(
      <AuthProvider>
        <TestConsumer />
      </AuthProvider>,
    );

    expect(screen.getByTestId("username")).toHaveTextContent("dummy");
    expect(screen.getByTestId("role")).toHaveTextContent("admin");
    expect(screen.getByTestId("id")).toHaveTextContent("1");
  });

  it("always returns dummy token", async () => {
    render(
      <AuthProvider>
        <TestConsumer />
      </AuthProvider>,
    );

    const btn = screen.getByText("Get Token");
    await act(async () => {
      fireEvent.click(btn);
    });

    expect(screen.getByTestId("token")).toHaveTextContent(
      "dummy-token-for-phase-1",
    );
  });

  it("always reports user as logged in", async () => {
    render(
      <AuthProvider>
        <TestConsumer />
      </AuthProvider>,
    );

    const btn = screen.getByText("Check Logged In");
    await act(async () => {
      fireEvent.click(btn);
    });

    expect(screen.getByTestId("loggedIn")).toHaveTextContent("true");
  });
});
